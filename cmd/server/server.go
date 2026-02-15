package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	commonMiddleware "github.com/hibare/GoCommon/v2/pkg/http/middleware"
	"github.com/hibare/GoGeoIP/cmd/server/handlers"
	"github.com/hibare/GoGeoIP/internal/config"
	"github.com/hibare/GoGeoIP/internal/maxmind"
	"github.com/spf13/cobra"
)

const (
	serverWriteTimeout    = 15 * time.Second
	serverReadTimeout     = 15 * time.Second
	serverIdleTimeout     = 60 * time.Second
	serverShutdownTimeout = 15 * time.Second
	middlewareTimeout     = 60 * time.Second
)

// Server represents the HTTP server.
type Server struct {
	cfg     *config.Config
	router  *chi.Mux
	ctx     context.Context
	maxmind *maxmind.Client
}

// NewServer creates a new Server instance.
func NewServer(ctx context.Context, cfg *config.Config, mm *maxmind.Client) *Server {
	return &Server{
		ctx:     ctx,
		cfg:     cfg,
		maxmind: mm,
	}
}

// Init initializes the server with handlers, routes and middleware.
func (s *Server) Init() error {
	// Initialize handlers
	geoIPHandler := handlers.NewGeoIP(s.maxmind)

	// Setup router
	s.router = chi.NewRouter()

	// Global middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(middlewareTimeout))
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middleware.CleanPath)
	s.router.Use(middleware.Heartbeat("/health"))

	// Register routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth)
		r.Group(func(r chi.Router) {
			r.Get("/ip", geoIPHandler.GetMyIP)
		})

		// Protected routes (with API key auth)
		r.Group(func(r chi.Router) {
			r.Use(func(h http.Handler) http.Handler {
				return commonMiddleware.TokenAuth(h, s.cfg.Server.APIKeys)
			})
			r.Get("/ip/{ip}", geoIPHandler.GetGeoIP)
		})
	})

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
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

		// Initialize MaxMind client
		mmClient := maxmind.NewClient(&config.Current.MaxMind, config.Current.Server.AssetDirPath)

		// Download DB if in production or missing
		if !config.Current.Server.IsDev {
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
		server := NewServer(ctx, config.Current, mmClient)
		if err := server.Init(); err != nil {
			return err
		}

		// Start serving
		return server.serve()
	},
}
