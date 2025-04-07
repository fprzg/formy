package services

import (
	"fmt"
	"strings"
	"time"
)

var allowedFieldTypes = map[string]bool{
	"int":      true,
	"float":    true,
	"string":   true,
	"datetime": true,
	"json":     true,
}

var simpleConstraints = map[string]bool{
	"required": true,
	"email":    true,
	"unique":   true,
	"hash":     true,
	"interval": true,
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is empty or whitespace", ErrInvalidInput)
	}

	return nil
}

func ValidateType(fieldType string) error {
	_, ok := allowedFieldTypes[fieldType]
	if !ok {
		return fmt.Errorf("invalid")
	}
	return nil
}

func ValidateConstraints(fieldType string, constraints []interface{}) error {
	for _, constraintAsInterface := range constraints {
		constraint, ok := constraintAsInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("constraint is not a map[string]interface{}: %w", ErrInvalidInput)
		}

		constraintName, ok := constraint["constraint_name"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid constraint_name: %w", ErrInvalidInput)
		}

		_, ok = simpleConstraints[constraintName]
		if !ok {
			return fmt.Errorf("unsupported constraint: %s: %w", constraintName, ErrInvalidInput)
		}

		if constraintName == "interval" {
			min, minOk := constraint["min"]
			max, maxOk := constraint["max"]

			if !minOk || !maxOk {
				return fmt.Errorf("interval constraint missing min or max: %w", ErrInvalidInput)
			}

			switch fieldType {
			case "int", "integer", "float":
				minF, minOk := min.(float64)
				maxF, maxOk := max.(float64)
				if !minOk || !maxOk {
					return fmt.Errorf("min/max should be numbers for field type %s: %w", fieldType, ErrInvalidInput)
				}
				if minF > maxF {
					return fmt.Errorf("min > max for interval in field type %s: %w", fieldType, ErrInvalidInput)
				}
			case "string":
				minF, minOk := min.(float64)
				maxF, maxOk := max.(float64)
				if !minOk || !maxOk {
					return fmt.Errorf("min/max should be numeric lengths for string field: %w", ErrInvalidInput)
				}
				if minF > maxF {
					return fmt.Errorf("min > max for string length interval: %w", ErrInvalidInput)
				}
			case "datetime":
				minStr, minOk := min.(string)
				maxStr, maxOk := max.(string)
				if !minOk || !maxOk {
					return fmt.Errorf("min/max should be strings for datetime interval: %w", ErrInvalidInput)
				}

				const layout = time.RFC3339
				minTime, err := time.Parse(layout, minStr)
				if err != nil {
					return fmt.Errorf("invalid min datetime: %w", ErrInvalidInput)
				}
				maxTime, err := time.Parse(layout, maxStr)
				if err != nil {
					return fmt.Errorf("invalid max datetime: %w", ErrInvalidInput)
				}
				if minTime.After(maxTime) {
					return fmt.Errorf("min datetime is after max datetime: %w", ErrInvalidInput)
				}
			case "bool":
				return fmt.Errorf("bool type cannot have intervals: %w", ErrInvalidInput)
			default:
				return fmt.Errorf("unsupported field type for interval: %s: %w", fieldType, ErrInvalidInput)
			}
		}
	}

	return nil
}
