package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	appErrors "github.com/hibare/GoGeoIP/cmd/server/errors"
	"github.com/hibare/GoGeoIP/cmd/server/middlewares"
	"github.com/hibare/GoGeoIP/cmd/server/utils"
	apikeys "github.com/hibare/GoGeoIP/internal/db/api_keys"
	"gorm.io/gorm"
)

type APIKeyHandler struct {
	db *gorm.DB
}

type CreateAPIKeyPayload struct {
	Name      string     `json:"name"                 validate:"required,min=1,max=255"`
	Scopes    []string   `json:"scopes,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type APIKeyCreateInput struct {
	Payload *CreateAPIKeyPayload `in:"body=json"`
}

type APIKeyIDInput struct {
	ID string `in:"path=id"`
}

func NewAPIKeyHandler(db *gorm.DB) *APIKeyHandler {
	return &APIKeyHandler{db: db}
}

// ListAPIKeys lists all API keys for the authenticated user.
func (h *APIKeyHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		http.Error(w, appErrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	params := r.URL.Query()
	params.Add("user_id", userID.String())

	keys, err := apikeys.ListAPIKeys(r.Context(), h.db, params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list API keys", "error", err)
		http.Error(w, appErrors.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, keys)
}

// CreateAPIKey creates a new API key.
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		http.Error(w, appErrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyCreateInput](r)
	if !ok {
		http.Error(w, appErrors.ErrReadingPayload.Error(), http.StatusBadRequest)
		return
	}

	apiKey := &apikeys.APIKey{
		UserID:    *userID,
		Name:      payload.Payload.Name,
		Scopes:    payload.Payload.Scopes,
		ExpiresAt: payload.Payload.ExpiresAt,
	}

	if err := apiKey.Validate(); err != nil {
		slog.ErrorContext(r.Context(), "invalid job payload", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, rawKey, err := apikeys.CreateAPIKey(r.Context(), h.db, apiKey)
	if err != nil {
		if errors.Is(err, apikeys.ErrDuplicateAPIKeyName) {
			http.Error(w, "API key name already exists", http.StatusConflict)
			return
		}
		slog.ErrorContext(r.Context(), "failed to create API key", "error", err)
		http.Error(w, appErrors.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, rawKey)
}

// RevokeAPIKey revokes an API key.
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		http.Error(w, appErrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyIDInput](r)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := apikeys.RevokeAPIKey(r.Context(), h.db, payload.ID, *userID)
	if err != nil {
		if errors.Is(err, apikeys.ErrAPIKeyNotFound) {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(r.Context(), "failed to revoke API key", "error", err)
		http.Error(w, appErrors.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAPIKey deletes an API key.
func (h *APIKeyHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		http.Error(w, appErrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyIDInput](r)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := apikeys.DeleteAPIKey(r.Context(), h.db, payload.ID, *userID)
	if err != nil {
		if errors.Is(err, apikeys.ErrAPIKeyNotFound) {
			http.Error(w, "API key not found", http.StatusNotFound)
			return
		}
		slog.ErrorContext(r.Context(), "failed to delete API key", "error", err)
		http.Error(w, appErrors.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
