package validator

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func Get() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})

	validate.RegisterValidation("cron", isValidCronExpression)
	return validate
}

func FirstError(err error) string {

	valErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ""
	}

	valErr := valErrors[0]

	switch valErr.Tag() {
	case "required":
		return fmt.Sprintf("%v is required", valErr.Field())
	case "email":
		return fmt.Sprintf("%v must be a valid email", valErr.Field())
	case "cron":
		return fmt.Sprintf("%v must be a valid cron expression", valErr.Field())
	default:
		return fmt.Sprintf("%v is not valid", valErr.Field())
	}
}

func isValidCronExpression(fl validator.FieldLevel) bool {

	expr := fl.Field().String()

	_, err := cron.ParseStandard(expr)
	return err == nil
}
