package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	commonErrors "github.com/hibare/GoCommon/v2/pkg/errors"
	commonHttp "github.com/hibare/GoCommon/v2/pkg/http"
	commonMiddleware "github.com/hibare/GoCommon/v2/pkg/http/middleware"
	commonIP "github.com/hibare/GoCommon/v2/pkg/net/ip"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/maxmind"
)

type Server struct{}

func NewServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Start() error {
	ctx := context.Background()
	router := chi.NewRouter()

	// Basic middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.NoCache)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(commonHttp.DefaultServerTimeout))
	router.Use(middleware.StripSlashes)
	router.Use(middleware.CleanPath)
	router.Use(middleware.Heartbeat(commonHttp.DefaultPingPath))

	// Use common security middleware
	router.Use(func(next http.Handler) http.Handler {
		return commonMiddleware.BasicSecurity(next, commonHttp.DefaultHTTPRequestSize)
	})

	// ToDo: @hibare Add favicon

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return commonMiddleware.TokenAuth(next, config.Current.Server.APIKeys)
		})
		r.Get("/ip", s.handlerMyIP)
		r.Get("/ip/{ip}", s.handleGeoIP)
	})

	srvAddr := fmt.Sprintf(":%d", config.Current.Server.ListenPort)
	srv := &http.Server{
		Handler:      router,
		Addr:         srvAddr,
		WriteTimeout: config.Current.Server.WriteTimeout,
		ReadTimeout:  config.Current.Server.ReadTimeout,
		IdleTimeout:  config.Current.Server.IdleTimeout,
	}

	slog.InfoContext(ctx, "GeoIP started", "address", srvAddr)

	// Schedule DB refresh
	if config.Current.DB.AutoUpdateEnabled {
		go maxmind.RunDBDownloadJob()
	}

	// Run server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(ctx, "failed to start server", "error", err)
			errChan <- err
		}
	}()

	// Check for startup errors
	select {
	case err := <-errChan:
		return err
	default:
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(ctx, commonHttp.DefaultServerShutdownGracePeriod)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Server shutdown failed", "error", err)
		return err
	}

	slog.InfoContext(ctx, "Server shutdown successfully")
	return nil
}

func (s Server) handlerMyIP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ipStr := r.RemoteAddr
	if !commonIP.IsPublicIP(ipStr) {
		commonHttp.WriteJsonResponse(w, http.StatusOK, maxmind.GeoIP{IP: ipStr, Remark: "Non-Public IP"})
		return
	}

	ipGeo, err := maxmind.IP2Geo(ipStr)

	if err != nil {
		slog.ErrorContext(ctx, "error fetching record", "ip", ipStr, "error", err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		return
	}
	commonHttp.WriteJsonResponse(w, http.StatusOK, ipGeo)
}

func (s Server) handleGeoIP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := chi.URLParam(r, "ip")

	ipGeo, err := maxmind.IP2Geo(ip)

	if err != nil {
		slog.ErrorContext(ctx, "error fetching record", "ip", ip, "error", err)
		if errors.Is(err, constants.ErrInvalidIP) {
			commonHttp.WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		commonHttp.WriteErrorResponse(w, http.StatusInternalServerError, commonErrors.ErrInternalServerError)
		return
	}
	commonHttp.WriteJsonResponse(w, http.StatusOK, ipGeo)
}
