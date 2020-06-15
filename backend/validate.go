package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"regexp"
)

const (
	nameRegexString string = "^[a-zA-Z]+(([',. -][a-zA-Z ])?[a-zA-Z]*)*$"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("name", isName); err != nil {
		log.Error(fmt.Sprintf("Error with registering custom validator: %v", err))
	}
}

func isName(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(nameRegexString)
	return reg.MatchString(fl.Field().String())
}

func (s *subscriber) Validate() error {
	if err := validate.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}
