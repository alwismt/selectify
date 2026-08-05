package httpx

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidationErrorsToMap(err error) map[string]string {
	errors := make(map[string]string)

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors
	}

	for _, fe := range ve {
		field := fe.Field()
		tag := fe.Tag()

		switch tag {
		case "required":
			errors[toSnakeCase(field)] = "This field is required"
		case "email":
			errors[toSnakeCase(field)] = "Invalid email address"
		case "min":
			errors[toSnakeCase(field)] = "Too short"
		case "max":
			errors[toSnakeCase(field)] = "Too long"
		default:
			errors[toSnakeCase(field)] = "Invalid value"
		}
	}

	return errors
}

func toSnakeCase(s string) string {
	var out []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}
		out = append(out, rune(strings.ToLower(string(r))[0]))
	}
	return string(out)
}
