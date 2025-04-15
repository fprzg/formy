package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"formy.fprzg.net/internal/controllers"
	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	"formy.fprzg.net/internal/types"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Port string
	Env  string

	m *models.Models
	s *services.Services
	c *controllers.Controllers
	e *echo.Echo
}

func NewServer(cfg types.AppConfig, db *sql.DB) (Server, error) {
	/*
		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", app.config.port),
			Handler:           app.routes(),
			ErrorLog:          log.New(app.logger, "", 0),
			IdleTimeout:       time.Minute,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 30 * time.Second,
		}
	*/

	var ctxDuration time.Duration
	if cfg.Env == "development" {
		ctxDuration = 1 * time.Hour
	} else {
		ctxDuration = 5 * time.Second
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JWTCustomClaims)
		},
		SigningKey:   []byte(cfg.JWTSecret),
		TokenLookup:  "cookie:jwt",
		ErrorHandler: errorHandler,
	}

	tm, err := services.NewTemplateManager(cfg.Env == "development", e)
	if err != nil {
		return Server{}, err
	}

	m, err := models.Get(db, e, ctxDuration)
	if err != nil {
		return Server{}, err
	}

	s, err := services.Get(cfg.JWTSecret, m, tm, e)
	if err != nil {
		return Server{}, err
	}

	c, err := controllers.Get(m, s, e, jwtConfig)
	if err != nil {
		return Server{}, err
	}

	err = insertDummyData(m)
	if err != nil {
		return Server{}, err
	}

	return Server{
		Port: cfg.Port,
		Env:  cfg.Env,
		e:    e,
		m:    m,
		s:    s,
		c:    c,
	}, nil
}

func errorHandler(c echo.Context, err error) error {
	return c.Redirect(http.StatusSeeOther, "/users/login")
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return nil
}

func (srv *Server) Serve() error {
	shutdownError := make(chan error)
	go srv.HandleSignals(shutdownError)

	srv.e.Logger.Info("starting server", map[string]string{
		"port": srv.Port,
		"env":  srv.Env,
	})

	err := srv.e.Start(srv.Port)
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	srv.e.Logger.Info("stopped server", map[string]string{
		"Port": srv.Port,
	})

	return nil

}

func (srv *Server) HandleSignals(shutdownError chan error) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	srv.e.Logger.Info("shutting down server", map[string]string{
		"signal": s.String(),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		shutdownError <- err
	}

	srv.e.Logger.Info("completing background tasks", map[string]string{
		"port": srv.Port,
	})

	//app.wg.Wait()
	shutdownError <- nil
}

func insertDummyData(m *models.Models) error {
	userID, err := models.InsertTestUser(m)
	if err != nil {
		return err
	}

	_, err = models.InsertTestForms(m, userID)
	if err != nil {
		return err
	}

	return nil
}
