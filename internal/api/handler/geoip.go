package handler

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/maxmind"
)

func GeoIP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	ipGeo, err := maxmind.IP2Geo(ip)

	if err != nil {
		log.Errorf("Error fetching record for ip %s, %s", ip, err)
		if errors.Is(err, constants.ErrInvalidIP) {
			render.Status(r, http.StatusBadRequest)
		} else {
			render.Status(r, http.StatusInternalServerError)
		}
		render.JSON(w, r, ErrorStruct{
			Message: err.Error(),
		})
		return
	}

	render.JSON(w, r, ipGeo)
}
