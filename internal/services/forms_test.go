package services

import (
	"log"
	"testing"

	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestFormServiceCreateForm(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	basicForm, err := types.JSONMapFromString(`
	{
		"name": "Basic Contact Form",
		"description": "Clients will use this form to contact us.",
		"fields": [
			{
				"field_name": "name",
				"field_type": "string",
				"field_constraints": [{ "constraint_name": "required"}] },
			{
				"field_name": "email",
				"field_type": "string",
				"field_constraints": [{"constraint_name": "unique"}, {"constraint_name": "required"}] },
			{
				"field_name": "phone_number",
				"field_type": "string",
				"field_constraints": [{ "constraint_name": "required"}] },
			{
				"field_name": "message",
				"field_type": "string",
				"field_constraints": []
			}
		]
	}
	`)

	if err != nil {
		log.Fatal(err.Error())
	}

	m := models.SetupTestDB(t)
	ms := GetModelServices(m)

	tests := []struct {
		TestName      string
		userID        int
		form          types.JSONMap
		expectedError error
	}{
		{
			TestName:      "Successful form upload",
			userID:        1,
			form:          basicForm,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err = ms.FormsServices.CreateForm(tt.userID, tt.form)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}
		})
	}
}


func submitHandle(c echo.Context) error {
	formValues := make(map[string]interface{})
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Failed to parse form data",
		})
	}