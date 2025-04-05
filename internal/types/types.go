package types

import "encoding/json"

type JSONMap map[string]interface{}

func JSONMapFromString(s string) (JSONMap, error) {
	var result JSONMap
	err := json.Unmarshal([]byte(s), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *JSONMap) ToJSONString() (string, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
