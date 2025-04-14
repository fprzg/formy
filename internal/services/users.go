package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UsersService interface {
}

type JWTCustomClaims struct {
	UserName string `json:"user_name"`
	UserID   int    `json:"user_id"`
	//Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (s *Services) UserLogin(c echo.Context) (string, error) {
	userName := c.FormValue("user_name")
	password := c.FormValue("password")

	userID, err := s.models.Users.Authenticate(userName, password)
	if err != nil {
		return "", echo.ErrUnauthorized
	}

	claims := &JWTCustomClaims{
		UserName: userName,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
