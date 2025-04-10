package controllers

import (
	"net/http"
	"strconv"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Controllers struct {
	models   *models.Models
	e        *echo.Echo
	services *services.Services
}

func New(m *models.Models) Controllers {
	c := Controllers{
		e:        echo.New(),
		models:   m,
		services: services.Get(m),
	}

	c.e.Use(middleware.Logger())
	c.e.Use(middleware.Recover())

	c.apiRoutes(c.e)
	c.frontendRoutes(c.e)
	c.uiRoutes(c.e)

	return c
}

func (c *Controllers) Start(cfg types.AppConfig) error {
	return c.e.Start(cfg.Port)
}

func (c *Controllers) pingHandle(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"status": "working",
	})
}

func (c *Controllers) apiRoutes(e *echo.Echo) {
	g := e.Group("/api")

	g.POST("/user/register", c.userRegisterHandle)
	g.PUT("/user/update", c.userUpdateHandle)

	g.POST("/form/create", c.formCreateHandle)
	g.POST("/form/get/:id", c.formGetHandle)
	g.PUT("/form/modify", c.formModifyHandle)

	g.POST("/submission/new/:id", c.submissionNewHandle)
}

func (c *Controllers) frontendRoutes(e *echo.Echo) {
	g := e.Group("/")

	g.GET("/", c.dummyFormHandler)
	g.GET("/ping", c.pingHandle)
}

func (c *Controllers) uiRoutes(e *echo.Echo) {
	_ = e.Group("/ui")

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
		//"form_type": formType,
		//"data": formValues,
	})
}

func (c *Controllers) dummyFormHandler(ctx echo.Context) error {
	html := `
	<!DOCTYPE html>
	<html>
	<body>
		<h2>Simple Form</h2>
		<form method="POST" action="/api/submit/1">
			<div>
				<label>Name:</label><br>
				<input type="text" name="name" required><br>
			</div>
			<div>
				<label>Email:</label><br>
				<input type="email" name="email" required><br>
			</div>
			<div>
				<label>Message:</label><br>
				<input type="text" name="message" required><br>
			</div>
			<div>
				<input type="submit" value="Submit">
			</div>
		</form>
	</body>
	</html>`

	return ctx.HTML(http.StatusOK, html)
}
