package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

// FormatValidationError translates validator error into user-friendly message
func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			tag := fieldError.Tag()
			param := fieldError.Param()

			message := formatFieldError(fieldName, tag, param, fieldError.Value())
			errors[strings.ToLower(fieldName)] = message
		}
	} else {
		// Fallback untuk non-validation errors
		errors["general"] = err.Error()
	}

	return errors
}

// formatFieldError creates user-friendly error message based on validation tag
func formatFieldError(field, tag, param string, value interface{}) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s tidak boleh kosong", humanizeFieldName(field))
	case "min":
		return fmt.Sprintf("%s minimal harus %s karakter", humanizeFieldName(field), param)
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", humanizeFieldName(field), param)
	case "email":
		return fmt.Sprintf("%s harus format email yang valid", humanizeFieldName(field))
	case "url":
		return fmt.Sprintf("%s harus URL yang valid", humanizeFieldName(field))
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", humanizeFieldName(field), param)
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", humanizeFieldName(field))
	case "uuid":
		return fmt.Sprintf("%s harus format UUID yang valid", humanizeFieldName(field))
	case "date":
		return fmt.Sprintf("%s harus format tanggal yang valid", humanizeFieldName(field))
	default:
		return fmt.Sprintf("%s tidak valid", humanizeFieldName(field))
	}
}

// humanizeFieldName converts field name to user-friendly name
func humanizeFieldName(field string) string {
	// Special cases for common abbreviations
	abbreviations := map[string]string{
		"ID":   "ID",
		"NIP":  "NIP",
		"NKKI": "NKKI",
	}
	
	if humanized, ok := abbreviations[field]; ok {
		return humanized
	}
	
	// Convert camelCase to Title Case with spaces
	result := ""
	for i, char := range field {
		if char >= 'A' && char <= 'Z' && i > 0 {
			result += " "
		}
		result += string(char)
	}
	return strings.ToLower(result)
}
