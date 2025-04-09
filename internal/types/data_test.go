package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldConstraint(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		TestName       string
		constraintJSON string
		expectedError  error
	}{
		{
			TestName:       "Email constraint",
			constraintJSON: `{"constraint_name": "email"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Unique constraint",
			constraintJSON: `{"constraint_name": "unique"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Required constraint",
			constraintJSON: `{"constraint_name": "required"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Integer interval constraint",
			constraintJSON: `{"constraint_name": "interval", "min": 0, "max": 150}`,
			expectedError:  nil,
		},
		{
			TestName:       "Dateime interval constraint",
			constraintJSON: `{"constraint_name": "interval", "min": "2022-01-01T00:00:00Z", "max": "2025-01-01T00:00:00Z"}`,
			expectedError:  nil,
		},
		{
			TestName: "String length interval constraint",
			//constraintJSON: `{"constraint_name": "interval", "min": "a", "max": "z"}`,
			constraintJSON: `{"constraint_name": "strlen", "min": 1, "max": 128}`,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			var fc FieldConstraint
			err := json.Unmarshal([]byte(tt.constraintJSON), &fc)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
			}
			fmt.Printf("%+v\n", fc)

		})
	}
}

func TestFieldData(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		TestName      string
		fieldContent  string
		expectedError error
	}{
		{
			TestName:      "Int",
			fieldContent:  `{"field_name": "Edad", "content": 30}`,
			expectedError: nil,
		},
		{
			TestName:      "Float",
			fieldContent:  `{"field_name": "Altura", "content": 1.75}`,
			expectedError: nil,
		},
		{
			TestName:      "String",
			fieldContent:  `{"field_name": "Nombre", "content": "Lucifer"}`,
			expectedError: nil,
		},
		{
			TestName:      "Datetime",
			fieldContent:  `{"field_name": "Nacimiento", "content": "2024-04-08T12:00:00Z"}`,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			var fc SubmissionField
			err := json.Unmarshal([]byte(tt.fieldContent), &fc)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
			}
			fmt.Printf("%+v\n", fc)
		})
	}
}
