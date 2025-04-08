package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type AppConfig struct {
	Port  string
	DBDir string
	Env   string
}

type UserData struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FormData struct {
	UserID      int         `json:"user_id"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"last_modified"`
	FormVersion int         `json:"form_version"`
	Fields      []FieldData `json:"fields"`
}

type FieldData struct {
	Name string `json:"field_name"`
	Type string `json:"field_type"`
	//Constraints string `json:"field_constraints"`
	Constraints []FieldConstraint
}

type FieldConstraint struct {
	Name string      `json:"constraint_name"`
	Min  interface{} `json:"min,omitempty"`
	Max  interface{} `json:"max,omitempty"`
}

func (fc *FieldConstraint) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Constraint name
	if v, ok := raw["constraint_name"]; ok {
		if err := json.Unmarshal(v, &fc.Name); err != nil {
			return err
		}
	}

	// Only try to unmarshal min/max if constraint is "interval"
	if fc.Name == "interval" {
		// Helper function to parse value
		parseValue := func(rawVal json.RawMessage) (interface{}, error) {
			// Try as int
			var i int
			if err := json.Unmarshal(rawVal, &i); err == nil {
				return i, nil
			}
			// Try as float
			var f float64
			if err := json.Unmarshal(rawVal, &f); err == nil {
				return f, nil
			}
			// Try as time
			var t time.Time
			if err := json.Unmarshal(rawVal, &t); err == nil {
				return t, nil
			}
			// Try as string
			var s string
			if err := json.Unmarshal(rawVal, &s); err == nil {
				return s, nil
			}
			return nil, fmt.Errorf("unknown type for value: %s", rawVal)
		}

		if rawMin, ok := raw["min"]; ok {
			v, err := parseValue(rawMin)
			if err != nil {
				return err
			}
			fc.Min = v
		}
		if rawMax, ok := raw["max"]; ok {
			v, err := parseValue(rawMax)
			if err != nil {
				return err
			}
			fc.Max = v
		}
	}

	return nil
}

type SubmissionData struct {
	FormID string
}
