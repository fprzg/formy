package types

import (
	"testing"
	"time"
)

func TestCheckvalueType(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		value          interface{}
		valueType      string
		expectedResult bool
	}{
		{
			name:           "int matches int",
			value:          123,
			valueType:      "int",
			expectedResult: true,
		},
		{
			name:           "int64 matches int64",
			value:          int64(123),
			valueType:      "int64",
			expectedResult: true,
		},
		{
			name:           "float64 matches float64",
			value:          3.14,
			valueType:      "float64",
			expectedResult: true,
		},
		{
			name:           "string matches string",
			value:          "hello",
			valueType:      "string",
			expectedResult: true,
		},
		{
			name:           "bool matches bool",
			value:          true,
			valueType:      "bool",
			expectedResult: true,
		},
		{
			name:           "time.Time matches time.Time",
			value:          now,
			valueType:      "time.Time",
			expectedResult: true,
		},
		{
			name:           "int does not match float64",
			value:          42,
			valueType:      "float64",
			expectedResult: false,
		},
		{
			name:           "nil value returns false",
			value:          nil,
			valueType:      "string",
			expectedResult: false,
		},
		{
			name:           "string does not match int",
			value:          "hello",
			valueType:      "int",
			expectedResult: false,
		},
		{
			name:           "time.Time does not match string",
			value:          now,
			valueType:      "string",
			expectedResult: false,
		},
		{
			name:           "unrecognized type name",
			value:          "hello",
			valueType:      "unknownType",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TypesMatch(tt.value, tt.valueType)
			if got != tt.expectedResult {
				t.Errorf("checkvalueType(%v, %q) = %v; want %v", tt.value, tt.valueType, got, tt.expectedResult)
			}
		})
	}
}
