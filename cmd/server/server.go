package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ggicci/httpin"
	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/hibare/GoGeoIP/cmd/server/handlers"
	"github.com/hibare/GoGeoIP/cmd/server/middlewares"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/constants"
	"github.com/hibare/GoGeoIP/internal/db"
	"github.com/hibare/GoGeoIP/internal/maxmind"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

//go:embed static/*
var staticFiles embed.FS

const (
	serverWriteTimeout    = 15 * time.Second
	serverReadTimeout     = 15 * time.Second
	serverIdleTimeout     = 60 * time.Second
	serverShutdownTimeout = 15 * time.Second
	middlewareTimeout     = 60 * time.Second
)

func getLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Server represents the HTTP server.
type Server struct {
	cfg     *config.Config
	router  *chi.Mux
	ctx     context.Context
	maxmind *maxmind.Client
	db      *gorm.DB
}

// NewServer creates a new Server instance.
func NewServer(ctx context.Context, cfg *config.Config, mm *maxmind.Client, db *gorm.DB) *Server {
	return &Server{
		ctx:     ctx,
		cfg:     cfg,
		maxmind: mm,
		db:      db,
	}
}

// Init initializes the server with handlers, routes and middleware.
func (s *Server) Init() error {
	apiKeyHandler := handlers.NewAPIKeyHandler(s.db)
	geoIPHandler := handlers.NewGeoIP(s.maxmind, s.cfg)
	authHandler, err := handlers.NewAuth(s.ctx, s.cfg, s.db)
	if err != nil {
		return fmt.Errorf("failed to create auth handler: %w", err)
	}

	s.router = chi.NewRouter()

	httpLogger := slog.Default()

	httpOptions := &httplog.Options{
		Level: getLogLevel(s.cfg.Logger.Level),
		Skip: func(req *http.Request, respStatus int) bool {
			return req.URL.Path == constants.HealthcheckPath
		},
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(httplog.RequestLogger(httpLogger, httpOptions))
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(middlewareTimeout))
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.Heartbeat(constants.HealthcheckPath))

	// Register routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public auth endpoints.
		r.Group(func(r chi.Router) {
			r.Get("/ip", geoIPHandler.GetMyIP)
			r.Route("/auth", func(r chi.Router) {
				r.With(httpin.NewInput(handlers.LoginInput{})).Get("/login", authHandler.Login)
				r.With(httpin.NewInput(handlers.CallbackInput{})).Get("/callback", authHandler.Callback)
				r.Post("/logout", authHandler.Logout)
				r.Get("/me", authHandler.Me)
			})
		})

		// Protected routes.
		r.Group(func(r chi.Router) {
			r.Use(middlewares.CookieAuthMiddleware)
			r.Get("/ip/{ip}", geoIPHandler.GetGeoIP)
			r.Get("/auth/me", authHandler.Me)

			// api keys routes
			r.Route("/api-keys", func(r chi.Router) {
				r.Get("/", apiKeyHandler.ListAPIKeys)
			})

			r.Route("/api-key", func(r chi.Router) {
				r.With(httpin.NewInput(handlers.APIKeyCreateInput{})).Post("/", apiKeyHandler.CreateAPIKey)
				r.With(httpin.NewInput(handlers.APIKeyIDInput{})).Post("/{id}/revoke", apiKeyHandler.RevokeAPIKey)
				r.With(httpin.NewInput(handlers.APIKeyIDInput{})).Delete("/{id}", apiKeyHandler.DeleteAPIKey)
			})
		})
	})

	// Serve static files only in production
	if s.cfg.Core.Environment == config.EnvironmentProduction {
		uiFS, err := fs.Sub(staticFiles, "static")
		if err != nil {
			log.Fatal(err)
		}
		fileServer := http.FileServer(http.FS(uiFS))
		s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			// We need to clean the path to ensure we look up the file correctly
			// and remove the leading slash because fs.Sub expects relative paths
			path := strings.TrimPrefix(r.URL.Path, "/")

			// Check if the file exists in the embedded filesystem
			if _, err := uiFS.Open(path); err == nil {
				// Case A: It is a real file (e.g., /assets/logo.png, /favicon.ico)
				// Serve it normally
				fileServer.ServeHTTP(w, r)
				return
			}

			// Case B: It is NOT a file (404 in FS).
			// This means it is likely an SPA route (e.g., /dashboard, /login).
			// We serve index.html so Vue Router can take over.

			// Force the path to root so FileServer serves index.html
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
		})
	} else {
		// In development, redirect all requests to UI dev server preserving path and query
		s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			targetURL := constants.UIAddress + r.URL.Path
			if r.URL.RawQuery != "" {
				targetURL += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, targetURL, http.StatusTemporaryRedirect)
		})
	}

	return nil
}

// serve starts the HTTP server with graceful shutdown.
func (s *Server) serve() error {
	addr := s.cfg.Server.GetAddr()

	srv := &http.Server{
		Handler:      s.router,
		Addr:         addr,
		WriteTimeout: serverWriteTimeout,
		ReadTimeout:  serverReadTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	slog.InfoContext(s.ctx, "Starting server", "address", addr)

	// Run our server in a goroutine so that it doesn't block.
	errChan := make(chan error, 1)
	go func() {
		var err error
		if s.cfg.Server.CertFile != "" && s.cfg.Server.KeyFile != "" {
			slog.InfoContext(s.ctx, "Starting server with TLS", "cert", s.cfg.Server.CertFile, "key", s.cfg.Server.KeyFile)
			err = srv.ListenAndServeTLS(s.cfg.Server.CertFile, s.cfg.Server.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			slog.ErrorContext(s.ctx, "failed to start server", "error", err)
			errChan <- err
		}
	}()

	// Check for immediate errors
	select {
	case err := <-errChan:
		return err
	default:
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(s.ctx, serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.ErrorContext(s.ctx, "Server shutdown failed", "error", err)
		return err
	}

	slog.InfoContext(s.ctx, "Server shutdown successfully")
	return nil
}

// ServeCmd represents the server command.
var ServeCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start API Server",
	Long:    "",
	Aliases: []string{"serve", "run"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		dbConn, err := db.New(ctx, config.Current)
		if err != nil {
			return err
		}

		// Initialize MaxMind client
		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Core.DataDir)

		// Download DB if in production or missing
		if config.Current.Core.Environment != config.EnvironmentDevelopment {
			if err := mmClient.DownloadAllDB(); err != nil {
				return err
			}
		} else {
			// Try to load existing DBs
			if err := mmClient.Load(); err != nil {
				slog.Warn("Failed to load MaxMind databases", "error", err)
			}
		}

		// Schedule background updates
		if config.Current.MaxMind.AutoUpdate {
			go mmClient.RunDBDownloadJob(ctx)
		}

		// Create and initialize server
		server := NewServer(ctx, config.Current, mmClient, dbConn.DB)
		if err := server.Init(); err != nil {
			return err
		}

		// Start serving
		return server.serve()
	},
}

func init() {
	httpin_integration.UseGochiURLParam("path", chi.URLParam)
}
