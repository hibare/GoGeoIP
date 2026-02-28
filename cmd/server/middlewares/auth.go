package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hibare/Waypoint/cmd/server/errors"
	"github.com/hibare/Waypoint/cmd/server/utils"
	"github.com/hibare/Waypoint/internal/auth"
	"github.com/hibare/Waypoint/internal/config"
	apikeys "github.com/hibare/Waypoint/internal/db/api_keys"
	"github.com/hibare/Waypoint/internal/db/users"
	"gorm.io/gorm"
)

// UserContextKey is the key for user information in request context.
type UserContextKey string

const UserKey UserContextKey = "user"

// GetAuthUser retrieves user claims from request context.
func GetAuthUser(r *http.Request) (*auth.UserJWTClaims, bool) {
	return utils.FromRequestContext[*auth.UserJWTClaims](r, UserKey)
}

// GetAuthUserID retrieves user ID from request context.
func GetAuthUserID(r *http.Request) (*uuid.UUID, bool) {
	user, ok := GetAuthUser(r)
	if !ok {
		return nil, false
	}
	userID := uuid.MustParse(user.UserID)
	return &userID, true
}

// UnifiedAuthMiddleware validates either API key or cookie authentication.
func UnifiedAuthMiddleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var claims *auth.UserJWTClaims

			// Try API key first
			if claims = tryAPIKeyAuth(ctx, db, r.Header.Get("Authorization")); claims != nil {
				goto authenticated
			}

			// Try cookie
			if claims = tryCookieAuth(r); claims != nil {
				goto authenticated
			}

			// Not authenticated
			http.Error(w, errors.ErrAuthenticationRequired.Error(), http.StatusUnauthorized)
			return

		authenticated:
			ctx = context.WithValue(ctx, UserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// tryCookieAuth attempts to authenticate via JWT cookie.
func tryCookieAuth(r *http.Request) *auth.UserJWTClaims {
	token, err := utils.GetJWTFromCookie(r)
	if err != nil {
		return nil
	}
	claims, err := auth.VerifyUserJWT(token)
	if err != nil {
		return nil
	}
	return claims
}

// tryAPIKeyAuth attempts to authenticate via API key in Authorization header.
func tryAPIKeyAuth(ctx context.Context, db *gorm.DB, authHeader string) *auth.UserJWTClaims {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil
	}

	apiKey := strings.TrimPrefix(authHeader, "Bearer ")
	if apiKey == "" {
		return nil
	}

	// Hash and lookup
	h := hmac.New(sha256.New, []byte(config.Current.Core.SecretKey))
	h.Write([]byte(apiKey))
	keyHash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	key, err := apikeys.GetAPIKeyByHash(ctx, db, keyHash)
	if err != nil {
		return nil
	}

	user, err := users.GetUserByID(ctx, db, key.UserID.String())
	if err != nil {
		return nil
	}

	// Update last used async (use background context to avoid cancellation)
	go func() { _ = apikeys.UpdateAPIKeyLastUsed(context.Background(), db, key.ID) }()

	return &auth.UserJWTClaims{
		UserID:     user.ID.String(),
		UserEmail:  user.Email,
		UserName:   strings.TrimSpace(user.FirstName + " " + user.LastName),
		UserGroups: user.Groups,
	}
}
