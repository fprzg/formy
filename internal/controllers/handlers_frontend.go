package controllers

import (
	"fmt"
	"net/http"
	"time"

	"formy.fprzg.net/internal/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// ///////////////////////////////////////////////
//
// # USERS HANDLERS
//
// ///////////////////////////////////////////////
func (ct *Controllers) handlerUsersRegisterPost(c echo.Context) error {
	return nil
}

func (ct *Controllers) handlerUsersRegisterGet(c echo.Context) error {
	td := services.NewTemplateData(c.Request())
	return ct.render(c, "users-register.tmpl.html", td)
}

func (ct *Controllers) handlerUsersLoginPost(c echo.Context) error {
	token, err := ct.services.UserLogin(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	expirationDate := time.Now().Add(time.Hour * 6).Unix()
	ct.setCookie(c, token, time.Unix(expirationDate, 0))

	return c.Redirect(http.StatusSeeOther, "/dash")
}

func (ct *Controllers) handlerUsersLoginGet(c echo.Context) error {
	td := services.NewTemplateData(c.Request())
	return ct.render(c, "users-login.tmpl.html", td)
}

func (ct *Controllers) handlerUsersLogout(c echo.Context) error {
	ct.setCookie(c, "", time.Unix(0, 0))

	return c.Redirect(http.StatusSeeOther, "/users/login")
}

// ///////////////////////////////////////////////
//
// # FRONTEND HANDLERS
//
// ///////////////////////////////////////////////
func (ct *Controllers) handlerHomePageGet(c echo.Context) error {
	td := services.NewTemplateData(c.Request())
	return ct.render(c, "home.tmpl.html", td)
}

func (ct *Controllers) handlerDashboardGet(c echo.Context) error {
	td := services.NewTemplateData(c.Request())

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*services.JWTCustomClaims)
	userData, err := ct.models.Users.Get(claims.UserID)
	if err != nil {
		return err
	}

	td.UserData = userData

	td.Dashboard = true
	return ct.render(c, "dash.tmpl.html", td)
}

// ///////////////////////////////////////////////
//
// # FORM HANDLERS
//
// ///////////////////////////////////////////////
func (ct *Controllers) handlerFormsCreatePost(c echo.Context) error {
	r := c.Request()
	formID, err := ct.services.ProcessForm(r, r.Context())
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("%d", formID))
}

func (ct *Controllers) formGetHandle(c echo.Context) error {
	return nil
}
