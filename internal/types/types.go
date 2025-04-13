package types

import "reflect"

func TypesMatch(value interface{}, valueType string) bool {
	if value == nil {
		return false
	}

	valueVal := reflect.ValueOf(value)
	valueKind := valueVal.Kind()
	valueTypeName := reflect.TypeOf(value).String()

	switch valueType {
	case "int":
		return valueKind == reflect.Int
	case "int64":
		return valueKind == reflect.Int64
	case "float64":
		return valueKind == reflect.Float64
	case "string":
		return valueKind == reflect.String
	case "bool":
		return valueKind == reflect.Bool
	case "time.Time":
		return valueTypeName == "time.Time"
	default:
		return valueTypeName == valueType
	}
}
