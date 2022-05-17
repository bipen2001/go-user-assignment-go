package utils

import (
	"regexp"

	"github.com/go-playground/validator"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[#$@!%&*?])[A-Za-z\d#$@!%&*?]{8,30}$/`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true

}
