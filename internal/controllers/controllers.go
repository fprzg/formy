package controllers

import (
	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	models    *models.Models
	services  *services.Services
	e         *echo.Echo
	public    *echo.Group
	protected *echo.Group
}

const StaticFilesDir = "../../public"

func Get(m *models.Models, s *services.Services, e *echo.Echo, jwtConfig echojwt.Config) (*Controllers, error) {
	c := &Controllers{
		models:    m,
		services:  s,
		e:         e,
		public:    e.Group(""),
		protected: e.Group("", echojwt.WithConfig(jwtConfig)),
	}

	c.staticFiles()
	c.apiRoutes()
	c.frontendRoutes()

	return c, nil
}

//
//
// ROUTES
//
//

func (c *Controllers) staticFiles() {
	c.public.Static("/static", StaticFilesDir)
}

func (c *Controllers) apiRoutes() {
	gPub := c.public.Group("/api")
	gProt := c.public.Group("/api")

	gPub.POST("/submissions/new/:id", c.handlerSubmissionsNew)
	gProt.POST("/ping", c.handlerPing)
}
func (c *Controllers) frontendRoutes() {
	gPub := c.public.Group("")
	gProt := c.protected.Group("")

	gPub.GET("/users/register", c.handlerUsersRegisterGet)
	gPub.POST("/users/register", c.handlerUsersRegisterPost)
	gPub.GET("/users/login", c.handlerUsersLoginGet)
	gPub.POST("/users/login", c.handlerUsersLoginPost)

	gProt.GET("/users/logout", c.handlerUsersLogout)
	gProt.POST("/users/logout", c.handlerUsersLogout)
	gProt.GET("/dash", c.handlerDashboard)
	gProt.POST("/form/create", c.handlerFormsCreatePost)
}
