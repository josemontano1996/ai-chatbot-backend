package utils

import "github.com/go-playground/validator/v10"

func ValidateStruct[T any](s T) error {
	validate := validator.New()
	err := validate.Struct(s)
	return err
}
