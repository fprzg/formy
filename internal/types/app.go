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
	Fields      []FormField `json:"fields"`
}

func (fd *FormData) GetFieldIndex(fieldName string) int {
	for idx, fieldDesc := range fd.Fields {
		if fieldDesc.Name == fieldName {
			return idx
		}
	}

	return -1
}

type FormField struct {
	Name        string            `json:"field_name"`
	Type        string            `json:"field_type"`
	Constraints []FieldConstraint `json:"field_constraints"`
}

type FieldConstraint struct {
	Name string      `json:"constraint_name"`
	Min  interface{} `json:"min,omitempty"`
	Max  interface{} `json:"max,omitempty"`
}

/*
func (fd FieldConstraint) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string      `json:"field_name"`
		Min  interface{} `json:"min,omitempty"`
		Max  interface{} `json:"max,omitempty"`
	}{
		Name: fd.Name,
		Min:  fd.Min,
		Max:  fd.Max,
	})
}
*/

func (fc *FieldConstraint) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if v, ok := raw["constraint_name"]; ok {
		if err := json.Unmarshal(v, &fc.Name); err != nil {
			return err
		}
	}

	if fc.Name == "interval" {
		parseValue := func(rawVal json.RawMessage) (interface{}, error) {
			var i int
			if err := json.Unmarshal(rawVal, &i); err == nil {
				return i, nil
			}

			var f float64
			if err := json.Unmarshal(rawVal, &f); err == nil {
				return f, nil
			}

			var t time.Time
			if err := json.Unmarshal(rawVal, &t); err == nil {
				return t, nil
			}

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
	} else if fc.Name == "strlen" {
		parseValue := func(rawVal json.RawMessage) (interface{}, error) {
			var i int
			if err := json.Unmarshal(rawVal, &i); err == nil {
				return i, nil
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
	ID             int               `json:"id"`
	FormID         int               `json:"form_id"`
	FormInstanceID int               `json:"form_instance_id"`
	Metadata       string            `json:"metadata"`
	SubmittedAt    string            `json:"submitted_at"`
	Fields         []SubmissionField `json:"fields"`
}

type SubmissionField struct {
	Name            string      `json:"field_name"`
	Type            string      `json:"field_type"`
	Content         interface{} `json:"field_content"`
	ContentAsString string      `json:"content_as_string"`
	Hash            string
	Unique          bool
}

/*
func (fd SubmissionData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name    string      `json:"field_name"`
		Content interface{} `json:"content"`
	}{
		Name:    fd.Name,
		Content: fd.Content,
	})
}
*/

func (fd *SubmissionField) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if nameRaw, ok := raw["field_name"]; ok {
		if err := json.Unmarshal(nameRaw, &fd.Name); err != nil {
			return err
		}
	}

	if contentRaw, ok := raw["content"]; ok {
		var i int
		if err := json.Unmarshal(contentRaw, &i); err == nil {
			fd.Content = i
			return nil
		}

		var f float64
		if err := json.Unmarshal(contentRaw, &f); err == nil {
			fd.Content = f
			return nil
		}

		var t time.Time
		if err := json.Unmarshal(contentRaw, &t); err == nil {
			fd.Content = t
			return nil
		}

		var s string
		if err := json.Unmarshal(contentRaw, &s); err == nil {
			fd.Content = s
			return nil
		}

		var v interface{}
		if err := json.Unmarshal(contentRaw, &v); err == nil {
			fd.Content = v
			return nil
		}

		return fmt.Errorf("unable to parse content")
	}

	return nil
}
