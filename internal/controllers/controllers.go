package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Form struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetRouter(env string) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/ping", pingHandle)

	e.POST("/api/user/register", userRegisterHandle)
	e.PUT("/api/user/update", userUpdateHandle)

	e.POST("/api/form/create", formCreateHandle)
	e.PUT("/api/form/modify", formModifyHandle)

	e.POST("/api/submit/:id", submitHandle)

	if env == "development" {
		e.GET("/", dummyFormHandler)
	}

	return e
}

func pingHandle(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"status": "working",
	})
}

func userRegisterHandle(c echo.Context) error {
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
}

func userUpdateHandle(c echo.Context) error {
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
}

func formCreateHandle(c echo.Context) error {
	form := new(Form)
	if err := c.Bind(form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid form data",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Form created successfully",
		"form":    form,
	})
}

func formModifyHandle(c echo.Context) error {
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
}

func submitHandle(c echo.Context) error {
	formValues := make(map[string]interface{})
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to parse form data",
		})
	}

	// Populate form values
	for key, values := range c.Request().PostForm {
		// If single value, store it directly; if multiple, store as slice
		if len(values) == 1 {
			formValues[key] = values[0]
		} else {
			formValues[key] = values
		}
	}

	/*
		// Get form type to determine how to process it
		formType, ok := formValues["form_type"].(string)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Form type is required",
			})
		}

			// Process based on form type
			switch formType {
			case "contact":
				// Validate contact form fields
				if _, ok := formValues["name"]; !ok {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"error": "Name is required for contact form",
					})
				}
				if _, ok := formValues["email"]; !ok {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"error": "Email is required for contact form",
					})
				}
				// Add more contact-specific validation/processing here

			case "feedback":
				// Validate feedback form fields
				if _, ok := formValues["rating"]; !ok {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"error": "Rating is required for feedback form",
					})
				}
				if _, ok := formValues["comment"]; !ok {
					return c.JSON(http.StatusBadRequest, echo.Map{
						"error": "Comment is required for feedback form",
					})
				}
				// Add more feedback-specific validation/processing here

			default:
				return c.JSON(http.StatusBadRequest, echo.Map{
					"error": "Unknown form type",
				})
			}
	*/

	// Return success response
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
		//"form_type": formType,
		//"data": formValues,
	})
}

func dummyFormHandler(c echo.Context) error {
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

	return c.HTML(http.StatusOK, html)
}
