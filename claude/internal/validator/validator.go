package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is the singleton validator instance.
var Validate *validator.Validate

// phone patterns for NG and US formats.
var (
	ngPhoneRegex = regexp.MustCompile(`^(\+?234|0)[789]\d{9}$`)
	usPhoneRegex = regexp.MustCompile(`^(\+?1)?[2-9]\d{2}[2-9]\d{6}$`)
)

func init() {
	Validate = validator.New()

	// Register custom phone validator for NG/US formats.
	Validate.RegisterValidation("phone_ng_us", func(fl validator.FieldLevel) bool {
		phone := strings.ReplaceAll(fl.Field().String(), " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		phone = strings.ReplaceAll(phone, "(", "")
		phone = strings.ReplaceAll(phone, ")", "")
		return ngPhoneRegex.MatchString(phone) || usPhoneRegex.MatchString(phone)
	})
}

// ValidateStruct validates a struct and returns field-level errors.
func ValidateStruct(s interface{}) map[string][]string {
	err := Validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string][]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := toSnakeCase(e.Field())
		var msg string
		switch e.Tag() {
		case "required":
			msg = "The " + field + " field is required."
		case "email":
			msg = "The " + field + " must be a valid email address."
		case "min":
			msg = "The " + field + " must be at least " + e.Param() + " characters."
		case "oneof":
			msg = "The selected " + field + " is invalid."
		case "eqfield":
			msg = "The " + field + " confirmation does not match."
		case "phone_ng_us":
			msg = "The " + field + " must be a valid NG or US phone number."
		default:
			msg = "The " + field + " field is invalid."
		}
		errors[field] = append(errors[field], msg)
	}
	return errors
}

// toSnakeCase converts PascalCase/camelCase to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
