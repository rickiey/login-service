package helper

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func IsValid(m interface{}) error {
	err := validate.Struct(m)
	if err != nil {
		return err
	}
	return nil
}
