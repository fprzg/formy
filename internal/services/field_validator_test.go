package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldValidatorValidateName(t *testing.T) {
	if testing.Short() {
		t.Skip("services: skiping whatever")
	}

	tests := []struct {
		TestName      string
		fieldName     string
		expectedError error
	}{
		{
			TestName:      "Valid field name",
			fieldName:     "valid_field_name",
			expectedError: nil,
		},
		{
			TestName:      "Invalid field name",
			fieldName:     "",
			expectedError: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			err := ValidateName(tt.fieldName)
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, ErrInvalidInput))
			}
		})
	}
}

func TestFieldValidatorValidateType(t *testing.T) {
	if testing.Short() {
		t.Skip("services: skiping whatever")
	}

	tests := []struct {
		TestName      string
		fieldType     string
		expectedError error
	}{
		{
			TestName: "Valid field type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
		})
	}
}

func TestFieldValidatorValidateConstraint(t *testing.T) {
	if testing.Short() {
		t.Skip("services: skiping whatever")
	}

	tests := []struct {
		TestName      string
		fieldType     string
		expectedError error
	}{
		{
			TestName: "Valid field type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
		})
	}
}
