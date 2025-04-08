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
			TestName:       "Valid ID",
			constraintJSON: `{"constraint_name": "email"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Invalid ID",
			constraintJSON: `{"constraint_name": "unique"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Valid ID",
			constraintJSON: `{"constraint_name": "required"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Valid ID",
			constraintJSON: `{"constraint_name": "interval", "min": 0, "max": 150}`,
			expectedError:  nil,
		},
		{
			TestName:       "Valid ID",
			constraintJSON: `{"constraint_name": "interval", "min": "2022-01-01T00:00:00Z", "max": "2025-01-01T00:00:00Z"}`,
			expectedError:  nil,
		},
		{
			TestName:       "Valid ID",
			constraintJSON: `{"constraint_name": "interval", "min": "a", "max": "z"}`,
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
