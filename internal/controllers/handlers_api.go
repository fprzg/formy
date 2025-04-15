package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (c *Controllers) handlerPingGet(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"status": "working",
	})
}

func (c *Controllers) handlerSubmissionsNewPost(ctx echo.Context) error {
	formID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	r := ctx.Request()
	submissionID, err := c.services.ProcessSubmission(formID, r, r.Context())
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message":       "success",
		"submission_id": submissionID,
	})
}
