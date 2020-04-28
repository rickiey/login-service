package helper

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func IsValid(m interface{}) error {
	err := validate.Struct(m)
	if err != nil {
		return err
	}
	return nil
}
