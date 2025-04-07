package controllers

import (
	"formy.fprzg.net/internal/services"
	"formy.fprzg.net/internal/types"
	"github.com/labstack/echo/v4"
)

var s services.ModelServices

func formsInsert(c echo.Context) {
	formAsJSON, err := types.JSONMapFromRequest(c.Request())
	if err != nil {
		return
	}
	s.FormsServices.CreateForm(formAsJSON)
}
