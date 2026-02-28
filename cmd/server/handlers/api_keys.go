package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	appErrors "github.com/hibare/Waypoint/cmd/server/errors"
	"github.com/hibare/Waypoint/cmd/server/middlewares"
	"github.com/hibare/Waypoint/cmd/server/utils"
	apikeys "github.com/hibare/Waypoint/internal/db/api_keys"
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
		commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	params := r.URL.Query()
	params.Add("user_id", userID.String())

	keys, err := apikeys.ListAPIKeys(r.Context(), h.db, params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list API keys", "error", err)
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, appErrors.ErrSomethingWentWrong)
		return
	}

	render.JSON(w, r, keys)
}

// CreateAPIKey creates a new API key.
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyCreateInput](r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, appErrors.ErrReadingPayload)
		return
	}

	apiKey := &apikeys.APIKey{
		UserID:    *userID,
		Name:      payload.Payload.Name,
		Scopes:    payload.Payload.Scopes,
		ExpiresAt: payload.Payload.ExpiresAt,
	}

	if err := apiKey.Validate(); err != nil {
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	_, rawKey, err := apikeys.CreateAPIKey(r.Context(), h.db, apiKey)
	if err != nil {
		if errors.Is(err, apikeys.ErrDuplicateAPIKeyName) {
			commonHttp.WriteErrorResponse(w, http.StatusConflict, err)
			return
		}
		slog.ErrorContext(r.Context(), "failed to create API key", "error", err)
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, appErrors.ErrSomethingWentWrong)
		return
	}

	render.JSON(w, r, rawKey)
}

// RevokeAPIKey revokes an API key.
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyIDInput](r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, appErrors.ErrReadingPayload)
		return
	}

	err := apikeys.RevokeAPIKey(r.Context(), h.db, payload.ID, *userID)
	if err != nil {
		if errors.Is(err, apikeys.ErrAPIKeyNotFound) {
			commonHttp.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}
		slog.ErrorContext(r.Context(), "failed to revoke API key", "error", err)
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, appErrors.ErrSomethingWentWrong)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAPIKey deletes an API key.
func (h *APIKeyHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetAuthUserID(r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusUnauthorized, appErrors.ErrUnauthorized)
		return
	}

	payload, ok := utils.InputFromContext[APIKeyIDInput](r)
	if !ok {
		commonHttp.WriteErrorResponse(w, http.StatusBadRequest, appErrors.ErrReadingPayload)
		return
	}

	err := apikeys.DeleteAPIKey(r.Context(), h.db, payload.ID, *userID)
	if err != nil {
		if errors.Is(err, apikeys.ErrAPIKeyNotFound) {
			commonHttp.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}
		slog.ErrorContext(r.Context(), "failed to delete API key", "error", err)
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, appErrors.ErrSomethingWentWrong)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
