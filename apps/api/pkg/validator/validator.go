package validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Decode decodes JSON from r.Body into dst and validates struct tags.
// Returns (validationFields, error). If validationFields is non-nil, respond with 422.
func Decode(r *http.Request, dst interface{}) (map[string]string, error) {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return nil, err
	}
	return Validate(dst)
}

// Validate runs struct validation and returns a field→message map on failure.
func Validate(dst interface{}) (map[string]string, error) {
	err := validate.Struct(dst)
	if err == nil {
		return nil, nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil, err
	}

	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		fields[field] = fieldMessage(fe)
	}
	return fields, nil
}

func fieldMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "value is too short (min " + fe.Param() + ")"
	case "max":
		return "value is too long (max " + fe.Param() + ")"
	case "uuid4":
		return "must be a valid UUID"
	case "oneof":
		return "must be one of: " + fe.Param()
	default:
		return "invalid value"
	}
}
