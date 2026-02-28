package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibare/Waypoint/cmd/server/errors"
	"github.com/hibare/Waypoint/cmd/server/utils"
	"github.com/hibare/Waypoint/internal/auth"
)

// UserContextKey is the key for user information in request context.
type UserContextKey string

const UserKey UserContextKey = "user"

// CookieAuthMiddleware validates JWT tokens from cookies and adds user info to context.
func CookieAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.GetJWTFromCookie(r)
		if err != nil {
			http.Error(w, errors.ErrAuthenticationRequired.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := auth.VerifyUserJWT(tokenString)
		if err != nil {
			utils.ClearAuthCookie(w) // Clear invalid cookie
			http.Error(w, errors.ErrInvalidAuthToken.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), UserKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthUser(r *http.Request) (*auth.UserJWTClaims, bool) {
	return utils.FromRequestContext[*auth.UserJWTClaims](r, UserKey)
}

func GetAuthUserID(r *http.Request) (*uuid.UUID, bool) {
	user, ok := GetAuthUser(r)
	if !ok {
		return nil, false
	}
	userID := uuid.MustParse(user.UserID)
	return &userID, true
}
