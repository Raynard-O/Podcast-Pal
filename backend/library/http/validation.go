package http

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func ValidationErrorResponse(c echo.Context, fe []validator.FieldError) error {
	var errors []ValidationError

	for _, err := range fe {
		var message string
		field := strings.ToLower(err.Field())

		switch err.ActualTag() {
		case "email":
			message = "is invalid"
		case "required":
			message = "is required"
		}

		ve := new(ValidationError)
		ve.Field = field
		ve.Message = field + " " + message

		errors = append(errors, *ve)
	}

	response := ErrorResponse{Errors: errors}
	return c.JSONPretty(http.StatusBadRequest, response, " ")
}
