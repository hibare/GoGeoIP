package apikeys

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/db"
	"gorm.io/gorm"
)

const (
	tableNameAPIKeys = "api_keys"
	maxNameLength    = 255
)

var (
	// ErrInvalidAPIKeyName is returned when the API key name is invalid.
	ErrInvalidAPIKeyName = errors.New("invalid API key name")

	// ErrAPIKeyNotFound is returned when the API key is not found.
	ErrAPIKeyNotFound = errors.New("API key not found")

	// ErrAPIKeyRevoked is returned when the API key is revoked.
	ErrAPIKeyRevoked = errors.New("API key is revoked")

	// ErrAPIKeyExpired is returned when the API key is expired.
	ErrAPIKeyExpired = errors.New("API key is expired")

	// ErrDuplicateAPIKeyName is returned when the API key name already exists for the user.
	ErrDuplicateAPIKeyName = errors.New("API key name already exists for this user")
)

type APIKeyStatus string

const (
	StatusActive  APIKeyStatus = "active"
	StatusRevoked APIKeyStatus = "revoked"
	StatusExpired APIKeyStatus = "expired"
)

type APIKey struct {
	ID         uuid.UUID  `json:"id"           gorm:"column:id;type:uuid;primaryKey"`
	UserID     uuid.UUID  `json:"user_id"      gorm:"column:user_id;type:uuid;not null"`
	Name       string     `json:"name"         gorm:"column:name;type:varchar(255);not null"`
	KeyHash    string     `json:"-"            gorm:"column:key_hash;type:varchar(255);unique;not null"`
	Scopes     []string   `json:"scopes"       gorm:"column:scopes;type:text[]"`
	ExpiresAt  *time.Time `json:"expires_at"   gorm:"column:expires_at;type:timestamp"`
	LastUsedAt *time.Time `json:"last_used_at" gorm:"column:last_used_at;type:timestamp"`
	RevokedAt  *time.Time `json:"revoked_at"   gorm:"column:revoked_at;type:timestamp"`
	State      string     `json:"state"        gorm:"column:state;type:varchar(50);not null;default:'active'"`
	CreatedAt  time.Time  `json:"created_at"   gorm:"autoCreateTime;column:created_at;not null"`
	UpdatedAt  time.Time  `json:"updated_at"   gorm:"autoUpdateTime;column:updated_at;not null"`
}

func (a *APIKey) TableName() string {
	return tableNameAPIKeys
}

func (a *APIKey) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// GetStatus returns the current status of the API key.
func (a *APIKey) GetStatus() APIKeyStatus {
	// First check explicit state field
	switch a.State {
	case string(StatusActive):
		// Even if state is active, check for revocation or expiration
		if a.RevokedAt != nil {
			return StatusRevoked
		}
		if a.ExpiresAt != nil && time.Now().UTC().After(*a.ExpiresAt) {
			return StatusExpired
		}
		return StatusActive
	case string(StatusRevoked):
		return StatusRevoked
	case string(StatusExpired):
		return StatusExpired
	default:
		// Fallback to computed status if state is unknown
		if a.RevokedAt != nil {
			return StatusRevoked
		}
		if a.ExpiresAt != nil && time.Now().UTC().After(*a.ExpiresAt) {
			return StatusExpired
		}
		return StatusActive
	}
}

// IsActive returns true if the API key is active.
func (a *APIKey) IsActive() bool {
	return a.GetStatus() == StatusActive
}

// Validate validates the API key data.
func (a *APIKey) Validate() error {
	if strings.TrimSpace(a.Name) == "" || len(a.Name) > maxNameLength {
		return ErrInvalidAPIKeyName
	}

	return nil
}

// generateAPIKey generates a new random API key with axon-api- prefix.
func generateAPIKey() string {
	return fmt.Sprintf("%s-api-%s", constants.ProgramIdentifier, uuid.New())
}

// HashAPIKey creates a SHA-256 hash of the API key for storage.
func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// CreateAPIKey creates a new API key.
func CreateAPIKey(ctx context.Context, db *gorm.DB, apiKey *APIKey) (*APIKey, string, error) {
	if err := apiKey.Validate(); err != nil {
		return nil, "", err
	}

	rawKey := generateAPIKey()
	apiKey.KeyHash = hashAPIKey(rawKey)

	if err := db.WithContext(ctx).Create(apiKey).Error; err != nil {
		return nil, "", err
	}

	return apiKey, rawKey, nil
}

// ListAPIKeys lists all API keys for a user.
func ListAPIKeys(ctx context.Context, tx *gorm.DB, params url.Values) ([]APIKey, error) {
	var apiKeys []APIKey

	qb := db.NewQueryBuilder()
	qb.RegisterIntField("id")
	qb.RegisterStringField("user_id")
	qb.RegisterStringField("name")
	qb.RegisterStringField("expires_at")
	qb.RegisterStringField("last_used_at")
	qb.RegisterTimeField("revoked_at")
	qb.RegisterTimeField("created_at")
	qb.RegisterTimeField("updated_at")

	// Parse URL query parameters into QueryOptions
	opts, err := qb.ParseQueryParams(params)
	if err != nil {
		return nil, err
	}

	err = tx.WithContext(ctx).
		Scopes(qb.Scope(opts)).
		Order("created_at DESC").
		Find(&apiKeys).Error

	return apiKeys, err
}

// RevokeAPIKey revokes an API key.
func RevokeAPIKey(ctx context.Context, db *gorm.DB, id string, userID uuid.UUID) error {
	now := time.Now().UTC()
	updates := map[string]any{
		"revoked_at": now,
		"state":      string(StatusRevoked),
	}

	result := db.WithContext(ctx).
		Model(&APIKey{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrAPIKeyNotFound
	}

	return nil
}

// DeleteAPIKey deletes an API key.
func DeleteAPIKey(ctx context.Context, db *gorm.DB, id string, userID uuid.UUID) error {
	result := db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&APIKey{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrAPIKeyNotFound
	}

	return nil
}

// UpdateAPIKeyLastUsed updates the last used timestamp.
func UpdateAPIKeyLastUsed(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	now := time.Now().UTC()
	return db.WithContext(ctx).
		Model(&APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", now).Error
}

// UpdateAPIKeyState updates the state of an API key.
func UpdateAPIKeyState(ctx context.Context, db *gorm.DB, id uuid.UUID, userID uuid.UUID, state string) error {
	result := db.WithContext(ctx).
		Model(&APIKey{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("state", state)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrAPIKeyNotFound
	}

	return nil
}

// MarkExpiredAPIKeys marks all expired API keys as expired.
func MarkExpiredAPIKeys(ctx context.Context, db *gorm.DB) error {
	now := time.Now().UTC()

	result := db.WithContext(ctx).
		Model(&APIKey{}).
		Where("expires_at IS NOT NULL AND expires_at < ? AND state = ?", now, string(StatusActive)).
		Update("state", string(StatusExpired))

	return result.Error
}
