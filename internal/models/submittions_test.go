package models

import (
	"testing"

	"formy.fprzg.net/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestSubmittionInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test.")
	}

	//m, err := GetTestModels()
	_, err := GetTestModels()
	assert.NoError(t, err)

	tests := []struct {
		TestName      string
		formID        int
		fields        []types.FieldData
		expectedError error
	}{
		{
			TestName: "Valid form insertion",
			formID:   1,
			fields: []types.FieldData{
				{
					Name:    "age",
					Content: 21,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			/*
				_, err := m.Submissions.Insert(tt.formID, "", tt.fields)
				if tt.expectedError == nil {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, tt.expectedError, err.Error())
				}
			*/
		})
	}
}
