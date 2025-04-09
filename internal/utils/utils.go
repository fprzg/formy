package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func PrintType(v interface{}) {
	fmt.Printf("%v", reflect.TypeOf(v))
}

func Print(v interface{}) {
	fmt.Printf("%v\n", v)
}

func ToJSON(v any) (*string, error) {
	asBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	asString := string(asBytes)
	return &asString, nil
}
