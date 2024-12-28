package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type Routers struct {
	log  *slog.Logger
	auth Auth
}

func NewRouters(log *slog.Logger, auth Auth) *Routers {

	return &Routers{
		log:  log,
		auth: auth,
	}
}

func (r *Routers) Login(c echo.Context) error {
	const op = "auth_service.http.Login"

	dto := new(LoginUserDto)

	if err := c.Bind(dto); err != nil {
		r.log.Info("Error bind dto", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request login auth_serivce" + err.Error()})
	}

	token, err := r.auth.Login(c.Request().Context(), dto.Email, dto.Password, dto.AppID)
	if err != nil {
		r.log.Info("Error login", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"Error login user" + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (r *Routers) Register(c echo.Context) error {
	const op = "auth_service.http.Register"

	dto := new(RegisterDto)

	if err := c.Bind(dto); err != nil {
		r.log.Info("Error bind dto", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request reginster auth_service" + err.Error()})
	}

	userID, err := r.auth.RegisterNewUser(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		r.log.Info("Error register user", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"Error register user" + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int64{"userID": userID})
}
