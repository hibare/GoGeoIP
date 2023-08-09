package handler

import (
	"errors"
	"net/http"

	commonErrors "github.com/hibare/GoCommon/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/pkg/http"
	commonIP "github.com/hibare/GoCommon/pkg/net/ip"

	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/maxmind"
	log "github.com/sirupsen/logrus"
)

func MyIP(w http.ResponseWriter, r *http.Request) {
	ipStr := r.RemoteAddr
	if !commonIP.IsPublicIP(ipStr) {
		commonHttp.WriteJsonResponse(w, http.StatusOK, maxmind.GeoIP{IP: ipStr, Remark: "Non-Public IP"})
		return
	}

	ipGeo, err := maxmind.IP2Geo(ipStr)

	if err != nil {
		log.Errorf("Error fetching record for ip %s, %s", ipStr, err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
		} else {
			commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		}
		return
	}
	commonHttp.WriteJsonResponse(w, http.StatusOK, ipGeo)
}
