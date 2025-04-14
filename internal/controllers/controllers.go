package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	models    *models.Models
	services  *services.Services
	e         *echo.Echo
	jwtSecret string
}

const StaticFilesDir = "../../public"

var secret string

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (c *Controllers) login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	if username != "jon" || password != "xx" {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		Name:  "Jon Snow",
		Admin: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(c.jwtSecret))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name

	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func Get(m *models.Models, s *services.Services, e *echo.Echo, jwtSecret string) (*Controllers, error) {
	c := &Controllers{
		models:    m,
		services:  s,
		e:         e,
		jwtSecret: jwtSecret,
	}

	e.POST("/login", c.login)
	e.GET("/", accessible)
	{

		r := e.Group("/restricted")
		config := echojwt.Config{
			//KeyFunc: getKey,
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(jwtCustomClaims)
			},
			SigningKey: []byte(jwtSecret),
		}
		r.Use(echojwt.WithConfig(config))
		r.GET("", restricted)
	}

	c.staticFiles()
	c.apiRoutes()
	c.frontendRoutes()
	c.uiRoutes()

	return c, nil
}

func (c *Controllers) pingHandle(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"status": "working",
	})
}

func (c *Controllers) staticFiles() {
	c.e.Static("/static", StaticFilesDir)
}

func (c *Controllers) apiRoutes() {
	g := c.e.Group("/api")

	g.POST("/user/register", c.userRegisterHandle)
	g.PUT("/user/update", c.userUpdateHandle)

	g.POST("/form/create", c.formCreateHandle)
	g.POST("/form/get/:id", c.formGetHandle)
	g.PUT("/form/modify", c.formModifyHandle)

	g.POST("/submission/new/:id", c.submissionNewHandle)
}

func (c *Controllers) frontendRoutes() {
	g := c.e.Group("")

	g.GET("/dash", c.dashboardHandler)
	g.GET("/ping", c.pingHandle)
}

func (c *Controllers) uiRoutes() {
	g := c.e.Group("/ui")

	g.POST("/form/new", c.uiFormNew)

	//e.GET("/user/:id", uiUsersGethandler)
	//e.GET("/form/:id", uiFormsGethandler)

}

//
//
// HANDLES
//
//

func (c *Controllers) userRegisterHandle(ctx echo.Context) error {
	/*
		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid user data",
			})
		}

		return c.JSON(http.StatusCreated, echo.Map{
			"message": "User registered successfully",
			"user":    user,
		})
	*/
	return nil
}

func (c *Controllers) userUpdateHandle(ctx echo.Context) error {
	/*
		user := new(User)
		if err := c.Bind(user); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid user data",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "User updated successfully",
			"user":    user,
		})
	*/
	return nil
}

func (c *Controllers) formCreateHandle(ctx echo.Context) error {
	r := ctx.Request()
	formID, err := c.services.Forms.ProcessForm(r, r.Context())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"form_id": formID,
	})
}

func (c *Controllers) formGetHandle(ctx echo.Context) error {
	return nil
}

func (c *Controllers) formModifyHandle(ctx echo.Context) error {
	/*
		form := new(Form)
		if err := c.Bind(form); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid form data",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Form modified successfully",
			"form":    form,
		})
	*/
	return nil
}

func (c *Controllers) submissionNewHandle(ctx echo.Context) error {
	formID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	r := ctx.Request()
	submissionID, err := c.services.Submissions.ProcessSubmission(formID, r, r.Context())
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message":       "success",
		"submission_id": submissionID,
	})
}

func (c *Controllers) dashboardHandler(ctx echo.Context) error {
	td := services.NewTemplateData(ctx.Request())
	td.Title = "Dashboard"
	td.Dashboard = true
	return c.render(ctx, "dash.tmpl.html", td)
}

func (c *Controllers) uiFormNew(ctx echo.Context) error {
	r := ctx.Request()
	formID, err := c.services.Forms.ProcessForm(r, r.Context())
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("%d", formID))
}

func (c *Controllers) render(ctx echo.Context, templateName string, td any) error {
	html, err := c.services.TemplateManager.ExecuteTemplate(templateName, td)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.HTML(http.StatusOK, html)
}
