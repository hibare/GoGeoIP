package middlewares

import (
	"net/http"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"

	"github.com/hibare/GoGeoIP/internal/api/handler"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/utils"
)

const AuthHeaderName = "Authorization"

func TokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Client: [%s] %s", r.RemoteAddr, r.RequestURI)

		apiKey := r.Header.Get(AuthHeaderName)

		if apiKey == "" {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, handler.ErrorStruct{
				Message: constants.ErrUnauthorized.Error(),
			})
			return
		}

		if utils.SliceContains(apiKey, config.Current.API.APIKeys) {
			next.ServeHTTP(w, r)
			return
		}

		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, handler.ErrorStruct{
			Message: constants.ErrUnauthorized.Error(),
		})
	})
}
