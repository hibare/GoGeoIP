package utils

import (
	"errors"
	"net/http"

	"github.com/hibare/Waypoint/internal/auth"
)

const (
	// OAuth2StateCookieName is the name of the cookie storing OAuth2 state.
	OAuth2StateCookieName = "oauth2_state"
	// OAuth2StateExpiration is the expiration time for OAuth2 state cookie in seconds.
	OAuth2StateExpiration = 300 // 5 minutes
)

var (
	ErrCookieNotFound = errors.New("cookie not found")
)

// SetCookie sets a cookie with the given name, value, and max age.
func SetCookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})
}

// GetCookie retrieves the value of a cookie by name.
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", ErrCookieNotFound
	}
	return cookie.Value, nil
}

// ClearCookie clears a cookie by setting its MaxAge to -1.
func ClearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // Expire immediately
	})
}

// SetAuthCookie sets the JWT token as a secure httpOnly cookie.
func SetAuthCookie(w http.ResponseWriter, tokenString string, maxAge int) {
	SetCookie(w, auth.JWTCookieName, tokenString, maxAge)
}

// ClearAuthCookie clears the authentication cookie.
func ClearAuthCookie(w http.ResponseWriter) {
	ClearCookie(w, auth.JWTCookieName)
}

// SetOAuth2StateCookie stores the OAuth2 state parameter in a cookie for verification.
func SetOAuth2StateCookie(w http.ResponseWriter, state string) {
	SetCookie(w, OAuth2StateCookieName, state, OAuth2StateExpiration)
}

// ClearOAuth2StateCookie clears the OAuth2 state cookie.
func ClearOAuth2StateCookie(w http.ResponseWriter) {
	ClearCookie(w, OAuth2StateCookieName)
}

// VerifyOAuth2State verifies the state parameter against the stored cookie.
func VerifyOAuth2State(r *http.Request, state string) error {
	value, err := GetCookie(r, OAuth2StateCookieName)
	if err != nil {
		return errors.New("OAuth2 state cookie not found")
	}

	if value != state {
		return errors.New("OAuth2 state parameter mismatch")
	}

	return nil
}

// GetJWTFromCookie extracts JWT token from the request cookie.
func GetJWTFromCookie(r *http.Request) (string, error) {
	return GetCookie(r, auth.JWTCookieName)
}
