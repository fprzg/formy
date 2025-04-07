package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/services"
	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Controllers struct {
	m *models.Models
	s *services.ModelServices
	e *echo.Echo
}

func GetControllers(m *models.Models) Controllers {
	c := Controllers{
		e: echo.New(),
		m: m,
		s: services.GetModelServices(m),
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
	g.PUT("/form/modify", c.formModifyHandle)

	g.POST("/submit/:id", c.submitHandle)
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

	if err := r.ParseForm(); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
	}

	formData := types.FormData{
		UserID:      r.FormValue("user_id"),
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
	}

	fieldNames := r.Form["field_name"]
	fieldTypes := r.Form["field_type"]
	fieldConstraints := r.Form["field_constraints"]

	if len(fieldNames) == 0 || len(fieldNames) != len(fieldTypes) || len(fieldNames) != len(fieldConstraints) {
		ctx.String(http.StatusBadRequest, "Invalid fields data")
		return nil
	}

	for i := range fieldNames {
		formData.Fields = append(formData.Fields, types.FieldData{
			Name:        fieldNames[i],
			Type:        fieldTypes[i],
			Constraints: fieldConstraints[i],
		})
	}

	userID, err := strconv.Atoi(formData.UserID)
	if err != nil {
		ctx.String(http.StatusBadRequest,
			fmt.Sprintf(`{ "status": "error", "message":  "%s" }`, err.Error()),
		)
		return err
	}

	formID, err := c.m.Forms.Insert(userID, formData.Name, formData.Description, formData.Fields)
	//formID, err := c.s.FormsServices.CreateForm(formData)
	if err != nil {
		ctx.String(http.StatusBadRequest,
			fmt.Sprintf(`{ "status": "error", "message":  "%s" }`, err.Error()),
		)
		return err
	}

	ctx.String(http.StatusOK,
		fmt.Sprintf(`{ "status": "OK", "form_id":  "%d" }`, formID),
	)

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

func (c *Controllers) submitHandle(ctx echo.Context) error {
	formValues, err := types.JSONMapFromRequest(ctx.Request())
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to parse form data",
		})
	}

	// Return success response
	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "success",
		//"form_type": formType,
		"data": formValues,
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
