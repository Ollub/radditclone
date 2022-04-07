package http_utils

import (
	"github.com/asaskevich/govalidator"
)

type Validator interface {
	IsValid() error
}

func Validate(in interface{}) error {
	_, err := govalidator.ValidateStruct(in)
	if err != nil {
		return err
	}

	validator, ok := in.(Validator)
	if !ok {
		return nil
	}
	err = validator.IsValid()
	if err != nil {
		return err
	}
	return nil
}
