package controllers

import (
	"net/http"
	"time"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	models    *models.Models
	services  *services.Services
	e         *echo.Echo
	JWTConfig echojwt.Config
	public    *echo.Group
	protected *echo.Group
}

const StaticFilesDir = "../../public"

func Get(m *models.Models, s *services.Services, e *echo.Echo, jwtConfig echojwt.Config) (*Controllers, error) {
	c := &Controllers{
		models:    m,
		services:  s,
		e:         e,
		JWTConfig: jwtConfig,
		protected: e.Group("", echojwt.WithConfig(jwtConfig)),
		public:    e.Group(""),
	}

	c.staticFiles()
	c.apiRoutes()
	c.frontendRoutes()

	return c, nil
}

func (ct *Controllers) render(c echo.Context, templateName string, td any) error {
	html, err := ct.services.TemplateManager.ExecuteTemplate(templateName, td)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.HTML(http.StatusOK, html)
}

func (ct *Controllers) setCookie(c echo.Context, value string, expirationDate time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = expirationDate

	c.SetCookie(cookie)
}

//
//
// ROUTES
//
//

func (c *Controllers) staticFiles() {
	pub := c.public.Group("/static")
	pub.Static("", StaticFilesDir)
}

func (c *Controllers) apiRoutes() {
	pub := c.public.Group("/api")
	pub.POST("/submissions/new/:id", c.handlerSubmissionsNewPost)

	prot := c.protected.Group("/api")
	prot.GET("/ping", c.handlerPingGet)
}
func (c *Controllers) frontendRoutes() {
	pub := c.public.Group("")
	pub.GET("", c.handlerHomePageGet)
	pub.GET("/users/register", c.handlerUsersRegisterGet)
	pub.POST("/users/register", c.handlerUsersRegisterPost)
	pub.GET("/users/login", c.handlerUsersLoginGet)
	pub.POST("/users/login", c.handlerUsersLoginPost)

	prot := c.protected.Group("")
	prot.GET("/users/logout", c.handlerUsersLogout)
	prot.POST("/users/logout", c.handlerUsersLogout)
	prot.GET("/dash", c.handlerDashboardGet)
	prot.POST("/form/create", c.handlerFormsCreatePost)
}
