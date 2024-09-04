package utils

import (
	"errors"
	"fmt"
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("first_is_a", func(fl validator.FieldLevel) bool {
			if value, ok := fl.Field().Interface().(string); ok {
				if value != "" && strings.Index(value, "a") == 0 {
					return true
				}
			}
			return false
		})
	}
}

func ParseValidationError(errs validator.ValidationErrors, target any) (errResult error) {
	fields := reflect.TypeOf(target).Elem()
	for _, fieldErr := range errs {
		field, _ := fields.FieldByName(fieldErr.Field())
		errMessageTag := fmt.Sprintf("%s_err", fieldErr.Tag())
		errMessage := field.Tag.Get(errMessageTag)

		if errMessage == "" {
			errMessage = field.Tag.Get("message")
		}

		if errMessage == "" {
			errMessage = fmt.Sprintf("%s:%s error", fieldErr.Field(), fieldErr.Tag())
		}
		fmt.Println("errMessage:", errMessage)
		errResult = response.AppendError(errResult, errors.New(errMessage))
	}
	return errResult
}
