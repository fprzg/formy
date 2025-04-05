package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormsInsertForm(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test.")
	}

	const contactFormFields = `
	[
		{ "field_name": "name", "field_type": "string", "constraints": ["required"] },
		{ "field_name": "email", "field_type": "string", "constraints": ["unique", "required"] },
		{ "field_name": "phone_number", "field_type": "string", "constraints": ["required"] },
		{ "field_name": "message", "field_type": "string", "constraints": [] }
	]
	`

	m := setupTestDB(t)

	tests := []struct {
		TestName      string
		userID        int
		name          string
		description   string
		fields        string
		expectedError error
	}{
		{
			TestName:      "Valid form insertion",
			userID:        1,
			name:          "Simple Contact Form",
			fields:        contactFormFields,
			expectedError: nil,
		},
		{
			TestName:      "Invalid userID",
			userID:        2,
			name:          "Simple Contact Form",
			fields:        contactFormFields,
			expectedError: ErrUserNotFound,
		},
		{
			TestName:      "Invalid name",
			userID:        1,
			name:          "",
			fields:        contactFormFields,
			expectedError: ErrInvalidInput,
		},
		{
			TestName:      "Invalid form fields",
			userID:        1,
			name:          "Simple Contact Form",
			fields:        "",
			expectedError: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Forms.InsertForm(tt.userID, tt.name, tt.description, tt.fields)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}
		})
	}
}

func TestFormsGetForm(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := setupTestDB(t)

	tests := []struct {
		TestName      string
		formID        int
		expectedError error
	}{
		{
			TestName:      "Valid ID",
			formID:        1,
			expectedError: nil,
		},
		{
			TestName:      "Invalid ID",
			formID:        0,
			expectedError: ErrFormNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			_, err := m.Forms.GetForm(tt.formID)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}

		})
	}
}

func TestFormsGetFormsByUser(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := setupTestDB(t)

	tests := []struct {
		TestName      string
		userID        int
		expectedError error
	}{
		{
			TestName:      "Valid ID",
			userID:        1,
			expectedError: nil,
		},
		{
			TestName:      "Empty forms",
			userID:        0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			_, err := m.Forms.GetFormsByUser(tt.userID)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}

		})
	}
}

func TestFormsUpdateFormName(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := setupTestDB(t)

	tests := []struct {
		TestName      string
		formID        int
		newName       string
		expectedError error
	}{
		{
			TestName:      "Successful name update",
			formID:        1,
			newName:       "New Form Name",
			expectedError: nil,
		},
		{
			TestName:      "Unsuccessful name update",
			formID:        0,
			newName:       "Another Form Name",
			expectedError: ErrFormNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Forms.UpdateFormName(tt.formID, tt.newName)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}

		})
	}
}

func TestFormsUpdateFormDescription(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := setupTestDB(t)

	tests := []struct {
		TestName       string
		formID         int
		newDescription string
		expectedError  error
	}{
		{
			TestName:       "Valid ID",
			formID:         1,
			newDescription: "Updated description",
			expectedError:  nil,
		},
		{
			TestName:       "Empty forms",
			formID:         0,
			newDescription: "Another description",
			expectedError:  ErrFormNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Forms.UpdateFormDescription(tt.formID, tt.newDescription)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}

		})
	}
}

func TestFormsUpdateFormDeleteForm(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	m := setupTestDB(t)

	tests := []struct {
		TestName      string
		formID        int
		expectedError error
	}{
		{
			TestName:      "Valid form deletion",
			formID:        1,
			expectedError: nil,
		},
		{
			TestName:      "Invalid form deletion",
			formID:        0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := m.Forms.DeleteForm(tt.formID)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, tt.expectedError, err.Error())
			}

		})
	}
}
