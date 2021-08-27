package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var passwordRegex = regexp.MustCompile(`^*[A-Z]*[a-z]*[0-9]*[_!@#$&*]`)

func ValidatePassword(f1 validator.FieldLevel) bool {
	val := false

	if len(f1.Field().String()) > 7 {
		val = passwordRegex.MatchString(f1.Field().String())
	}
	return val
}
