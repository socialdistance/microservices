package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	httprouters "lib_isod_v2/auth_service/internal/http"

	"github.com/arl/statsviz"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	m       *http.ServeMux
	log     *slog.Logger
	e       *echo.Echo
	routers *httprouters.Routers
	host    string
	port    string
}

func New(log *slog.Logger, host, port string, routers *httprouters.Routers) *Server {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info("request",
				slog.String("URI", v.URI),
				slog.Int("status", v.Status),
				slog.String("remote ip", v.RemoteIP),
			)

			return nil
		},
	}))

	mux := http.NewServeMux()
	statsviz.Register(mux)

	return &Server{
		m:       mux,
		log:     log,
		e:       e,
		routers: routers,
		host:    host,
		port:    port,
	}
}

// MustRun runs HTTP server and panics if any error occurs.
func (s *Server) MustRun() {
	const op = "http.Server.MustRun"

	s.log.Info(op, slog.String("Start", "server"))

	if err := s.Start(); err != nil {
		panic(err)
	}
}

func (s *Server) Start() error {
	const op = "http.Server.Start"

	if err := s.e.Start(fmt.Sprintf(":%s", s.port)); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s server stopped: %w", op, err)
	}

	return nil
}

func (s *Server) Stop() error {
	const op = "http.Server.Stop"

	optCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.log.Info("stopping", op, "http server")

	if err := s.e.Shutdown(optCtx); err != nil {
		return fmt.Errorf("%s could not shutdown server gracefuly: %w", op, err)
	}

	return nil
}

func (s *Server) BuildRouters() {
	debug := s.e.Group("/debug")
	debug.GET("/statsviz/", echo.WrapHandler(s.m))
	debug.GET("/statsviz/*", echo.WrapHandler(s.m))

	api := s.e.Group("/api/auth")

	api.POST("/login", s.routers.Login)
	api.POST("/register", s.routers.Register)
	api.GET("/apps", s.routers.ListApps)
}
