package types

import (
	"encoding/json"
	"net/http"
)

type JSONMap map[string]interface{}

func JSONHasField(json map[string]interface{}, field string) bool {
	_, ok := json[field]
	return ok
}

func JSONMapFromString(s string) (JSONMap, error) {
	var result JSONMap
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func JSONMapFromRequest(r *http.Request) (JSONMap, error) {
	formValues := make(JSONMap)

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	for key, values := range r.PostForm {
		// If single value, store it directly; if multiple, store as slice
		if len(values) == 1 {
			formValues[key] = values[0]
		} else {
			formValues[key] = values
		}
	}

	return formValues, nil
}

func (m *JSONMap) ToJSONString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
