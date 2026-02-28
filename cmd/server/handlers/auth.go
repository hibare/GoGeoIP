package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/render"
	apperrors "github.com/hibare/Waypoint/cmd/server/errors"
	"github.com/hibare/Waypoint/cmd/server/middlewares"
	"github.com/hibare/Waypoint/cmd/server/utils"
	"github.com/hibare/Waypoint/internal/auth"
	"github.com/hibare/Waypoint/internal/config"
	"github.com/hibare/Waypoint/internal/db/users"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	ErrAuthNotEnabled  = errors.New("OIDC authentication is not enabled")
	ErrStateMissing    = errors.New("state parameter missing")
	ErrInvalidState    = errors.New("invalid state parameter")
	ErrAuthCodeMissing = errors.New("authorization code not provided")
	ErrFailedExchange  = errors.New("failed to exchange code for token")
	ErrNoIDToken       = errors.New("no id_token field in oauth2 token")
	ErrFailedVerify    = errors.New("failed to verify ID token")
	ErrFailedClaims    = errors.New("failed to extract user claims")
	ErrNoAuthCode      = errors.New("no authorization code found")
	ErrMissingUserInfo = errors.New("missing user info")
)

const (
	defaultExpiration = 4 * time.Hour
	rootPath          = "/"
	err500route       = "/500"
)

// Auth represents the authentication handler.
type Auth struct {
	cfg          *config.Config
	db           *gorm.DB
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	provider     *oidc.Provider
}

type Claims struct {
	Sub    string   `json:"sub"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Groups []string `json:"groups,omitempty"` // OIDC groups claim
	Exp    int64    `json:"exp"`              // Expiration time from ID token
}

func (c *Claims) PostProcess() {
	if c.Exp == 0 {
		c.Exp = time.Now().UTC().Add(defaultExpiration).Unix()
	}
}

func (c *Claims) GetExpirationTime() time.Time {
	return time.Unix(c.Exp, 0)
}

// AuthResponse represents the authentication response.
type AuthResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Token   string      `json:"token,omitempty"`
	User    *users.User `json:"user,omitempty"`
}

// createOrUpdateUser creates a new user or updates an existing one based on email.
func (a *Auth) createOrUpdateUser(ctx context.Context, email, name string, groups []string) (*users.User, error) {
	if email == "" || name == "" {
		return nil, ErrMissingUserInfo
	}

	// Parse name into first and last name
	nameParts := strings.SplitN(name, " ", 2) //nolint:mnd // splitting into first and last name
	firstName := nameParts[0]
	lastName := ""
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	// Check if user exists
	existingUser, err := users.GetUserByEmail(ctx, a.db, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// User doesn't exist, create new user
		user := &users.User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Groups:    groups,
			LastLogin: time.Now().UTC(),
		}

		if err := users.CreateUser(ctx, a.db, user); err != nil {
			return nil, err
		}

		return user, nil
	} else {
		// User exists, update user
		updates := &users.User{
			FirstName: firstName,
			LastName:  lastName,
			Groups:    groups,
			LastLogin: time.Now().UTC(),
		}

		if err := users.UpdateUser(ctx, a.db, existingUser.ID.String(), updates); err != nil {
			return nil, err
		}

		// Return updated user
		return users.GetUserByID(ctx, a.db, existingUser.ID.String())
	}
}

func (a *Auth) redirect(w http.ResponseWriter, r *http.Request, redirect string, message string) {
	redirectURL := rootPath

	// If we have a redirect path from state, append it to the base redirect URL
	parsedURL, err := url.Parse(redirectURL)
	if err == nil {
		if redirect != "" {
			parsedURL.Path = redirect
		}

		if message != "" {
			q := parsedURL.Query()
			q.Set("message", message)
			parsedURL.RawQuery = q.Encode()
		}
	}

	redirectURL = parsedURL.String()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// LoginInput represents the input for login request.
type LoginInput struct {
	Redirect string `in:"query=redirect"`
}

// CallbackInput represents the input for OAuth callback request.
type CallbackInput struct {
	State string `in:"query=state"`
	Code  string `in:"query=code"`
}

// Login redirects to the OIDC provider for authentication.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get redirect parameter using httpin
	payload, ok := utils.InputFromContext[LoginInput](r)
	if !ok {
		a.redirect(w, r, err500route, apperrors.ErrReadingPayload.Error())
		return
	}

	redirect := payload.Redirect

	// Create state with redirect information
	state, err := utils.NewState()
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate state parameter", "error", err)
		a.redirect(w, r, err500route, apperrors.ErrSomethingWentWrong.Error())
		return
	}

	// Add redirect information to state
	state.AddInfo("redirect", redirect)

	// Encode state
	encodedState := state.String()

	// Store encoded state in cookie for verification in callback
	utils.SetOAuth2StateCookie(w, encodedState)

	// Use OAuth2 library's AuthCodeURL method with encoded state parameter
	authURL := a.oauth2Config.AuthCodeURL(encodedState)

	render.JSON(w, r, map[string]string{
		"redirect_url": authURL,
	})
}

// Callback handles the OIDC callback and exchanges code for token.
func (a *Auth) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get callback parameters using httpin
	payload, ok := utils.InputFromContext[CallbackInput](r)
	if !ok {
		a.redirect(w, r, err500route, apperrors.ErrReadingPayload.Error())
		return
	}

	encodedState := payload.State
	if encodedState == "" {
		a.redirect(w, r, err500route, ErrStateMissing.Error())
		return
	}

	// Decode state to get redirect information
	state, err := utils.NewStateFromEncode(encodedState)
	if err != nil {
		slog.ErrorContext(ctx, "failed to decode state parameter", "error", err)
		a.redirect(w, r, err500route, ErrInvalidState.Error())
		return
	}

	// Get redirect information from state
	redirect, _ := state.GetInfo("redirect")

	// Verify state parameter matches the one stored in cookie
	if verifyErr := utils.VerifyOAuth2State(r, encodedState); verifyErr != nil {
		slog.ErrorContext(ctx, "invalid state parameter", "error", verifyErr)
		a.redirect(w, r, err500route, ErrInvalidState.Error())
		return
	}

	// Clear the state cookie after successful verification
	utils.ClearOAuth2StateCookie(w)

	code := payload.Code
	if code == "" {
		a.redirect(w, r, err500route, ErrNoAuthCode.Error())
		return
	}

	// Exchange code for token and extract user claims
	claims, exchangeErr := a.exchangeCodeForClaims(ctx, payload.Code)
	if exchangeErr != nil {
		slog.ErrorContext(ctx, "failed to exchange code for claims", "error", exchangeErr)
		a.redirect(w, r, err500route, exchangeErr.Error())
		return
	}

	claims.PostProcess()

	// Create or update user in database
	user, err := a.createOrUpdateUser(ctx, claims.Email, claims.Name, claims.Groups)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create / update user", "error", err)
		a.redirect(w, r, err500route, apperrors.ErrSomethingWentWrong.Error())
		return
	}

	// Create JWT token with expiration matching the ID token
	jwtToken, err := auth.CreateUserJWT(user, claims.GetExpirationTime())
	if err != nil {
		slog.ErrorContext(ctx, "failed to create JWT token", "error", err)
		a.redirect(w, r, err500route, apperrors.ErrSomethingWentWrong.Error())
		return
	}

	utils.SetAuthCookie(w, jwtToken, int(time.Until(claims.GetExpirationTime()).Seconds()))
	a.redirect(w, r, redirect, "")
}

func (a *Auth) exchangeCodeForClaims(ctx context.Context, code string) (*Claims, error) {
	if code == "" {
		return nil, ErrNoAuthCode
	}

	token, err := a.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, ErrFailedExchange
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, ErrNoIDToken
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, ErrFailedVerify
	}

	var claims Claims
	if err = idToken.Claims(&claims); err != nil {
		return nil, ErrFailedClaims
	}

	// Always get additional user info from userinfo endpoint for complete profile data
	userInfo, err := a.getUserInfo(ctx, token.AccessToken)
	if err != nil {
		slog.WarnContext(ctx, "failed to get userinfo, using ID token claims only", "error", err)
	} else {
		// Merge userinfo with ID token claims, preferring userinfo for profile data
		if userInfo.Email != "" {
			claims.Email = userInfo.Email
		}
		if userInfo.Name != "" {
			claims.Name = userInfo.Name
		}
		if len(userInfo.Groups) > 0 {
			claims.Groups = userInfo.Groups
		}
	}

	claims.PostProcess()
	return &claims, nil
}

// getUserInfo fetches additional user information from the OIDC userinfo endpoint.
func (a *Auth) getUserInfo(ctx context.Context, accessToken string) (*Claims, error) {
	// Get userinfo from the provider
	userInfo, err := a.provider.UserInfo(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken}))
	if err != nil {
		return nil, err
	}

	var claims Claims
	if cErr := userInfo.Claims(&claims); cErr != nil {
		return nil, cErr
	}

	return &claims, nil
}

// Logout handles user logout.
func (a *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the authentication cookie
	utils.ClearAuthCookie(w)

	render.NoContent(w, r)
}

func (a *Auth) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middlewares.GetAuthUser(r)
	if !ok {
		http.Error(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}
	user, err := users.GetUserByID(r.Context(), a.db, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, apperrors.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, user)
}

// NewAuth creates a new authentication handler with initialized OIDC components.
func NewAuth(ctx context.Context, cfg *config.Config, db *gorm.DB) (*Auth, error) {
	auth := &Auth{cfg: cfg, db: db}

	// Initialize OIDC provider (automatically performs discovery)
	provider, err := oidc.NewProvider(ctx, cfg.OIDC.IssuerURL)
	if err != nil {
		return nil, err
	}

	// Create OAuth2 config using provider's discovered endpoints
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.OIDC.ClientID,
		ClientSecret: cfg.OIDC.ClientSecret,
		RedirectURL:  cfg.OIDC.GetRedirectURI(cfg.Server.BaseURL),
		Scopes:       cfg.OIDC.Scopes,
		Endpoint:     provider.Endpoint(), // Automatically discovered endpoints
	}

	// Create ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.OIDC.ClientID,
	})

	auth.oauth2Config = oauth2Config
	auth.verifier = verifier
	auth.provider = provider
	return auth, nil
}
