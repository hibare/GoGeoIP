package handler

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/maxmind"
)

func GeoIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	ipGeo, err := maxmind.IP2Geo(ip)

	if err != nil {
		log.Errorf("Error fetching record for ip %s, %s", ip, err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
		} else {
			commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		}
		return
	}
	commonHttp.WriteJsonResponse(w, http.StatusOK, ipGeo)
}
