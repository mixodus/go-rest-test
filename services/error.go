package services

import (
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Should be greater than " + fe.Param() + " characters"
	case "max":
		return "Should be less than " + fe.Param() + " characters"
	case "email":
		return "Should be a valid email address"
	case "number":
		return "Should be a number"
	case "file":
		return "Should be a file"
	}

	return "Unknown error"
}
