package shared

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

func Validate(payload any) error {
	if err := validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			field := strings.ToLower(validationErrors[0].Field())
			tag := validationErrors[0].Tag()
			if tag == "required" {
				return ValidationError(fmt.Sprintf("%s is required", field))
			}
			return ValidationError(fmt.Sprintf("%s is invalid", field))
		}
		return err
	}
	return nil
}
