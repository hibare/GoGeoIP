package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/maxmind"
)

// GeoIP handles GeoIP-related requests.
type GeoIP struct {
	maxmind *maxmind.Client
}

// NewGeoIP creates a new GeoIP handler.
func NewGeoIP(mm *maxmind.Client) *GeoIP {
	return &GeoIP{
		maxmind: mm,
	}
}

// GetGeoIP handles requests to get GeoIP information for a specific IP.
func (h *GeoIP) GetGeoIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	ipGeo, err := h.maxmind.IP2Geo(ip)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error fetching record for ip", "ip", ip, "error", err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
		} else {
			commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		}
		return
	}
	commonHttp.WriteJSONResponse(w, http.StatusOK, ipGeo)
}

// GetMyIP handles requests to get GeoIP information for the requester's IP.
func (h *GeoIP) GetMyIP(w http.ResponseWriter, r *http.Request) {
	ipStr := r.RemoteAddr

	ipGeo, err := h.maxmind.IP2Geo(ipStr)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error fetching record for ip", "ip", ipStr, "error", err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
		} else {
			commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		}
		return
	}
	commonHttp.WriteJSONResponse(w, http.StatusOK, ipGeo)
}
