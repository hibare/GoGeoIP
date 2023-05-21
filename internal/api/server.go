package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/hibare/GoGeoIP/internal/api/handler"
	"github.com/hibare/GoGeoIP/internal/api/middlewares"
	"github.com/hibare/GoGeoIP/internal/config"
)

type App struct {
	Router    *chi.Mux
	Validator *validator.Validate
}

func (a *App) setRouters() {
	a.Get("/api/v1/health", a.HealthCheck, false)
	a.Get("/api/v1/ip/{ip}", a.GeoIP, true)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request), protected bool) {
	if protected {
		pr := chi.NewRouter()
		pr.Use(middlewares.TokenAuth)
		pr.Get(path, f)
		a.Router.Mount("/", pr)
	} else {
		a.Router.Get(path, f)
	}
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request), protected bool) {
	a.Router.Post(path, f)
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request), protected bool) {
	a.Router.Put(path, f)
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request), protected bool) {
	a.Router.Delete(path, f)
}

func (a *App) HealthCheck(w http.ResponseWriter, r *http.Request) {
	handler.HealthCheck(w, r)
}

func (a *App) GeoIP(w http.ResponseWriter, r *http.Request) {
	handler.GeoIP(w, r)
}

// App initialize with predefined configuration
func (a *App) Initialize() {
	validator := validator.New()
	a.Validator = validator

	a.Router = chi.NewRouter()
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(middleware.Timeout(60 * time.Second))
	a.Router.Use(middleware.NoCache)
	a.Router.Use(middleware.StripSlashes)
	a.Router.Use(middleware.CleanPath)

	a.setRouters()
}

func (a *App) Serve() {
	wait := time.Second * 15
	addr := fmt.Sprintf("%s:%d", config.Current.API.ListenAddr, config.Current.API.ListenPort)

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}

	log.Infof("Listening for address %s on port %d", config.Current.API.ListenAddr, config.Current.API.ListenPort)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Info(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
}
