package tool

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

// ValidateErrorResult 自定义校验失败信息
type ValidateErrorResult struct {
	Error       bool
	FailedField string
	Value       interface{}
	Tag         string
}

var globalValidator *validator.Validate

func init() {
	globalValidator = validator.New()
	_ = globalValidator.RegisterValidation("valid-path", validPathValidation)
}

func validPathValidation(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return ValidArticlePath(s)
}

func ValidArticlePath(str string) bool {
	match, err := regexp.Match(`^[a-z0-9:@._-]+$`, []byte(str))
	if err != nil {
		return false
	}
	return match
}

// 校验字段
func validate(data interface{}) []ValidateErrorResult {
	var validationErrors []ValidateErrorResult
	errs := globalValidator.Struct(data)
	if errs != nil {
		// 收集所有校验失败的字段信息
		for _, err := range errs.(validator.ValidationErrors) {
			elem := ValidateErrorResult{
				Error:       true,
				FailedField: err.Field(), // 获取字段名
				Value:       err.Value(), // 获取字段值
				Tag:         err.Tag(),   // 获取错误信息
			}
			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}

// GetValidateError 校验并获取错误信息
func GetValidateError(body interface{}) error {
	if errorResults := validate(body); len(errorResults) > 0 && errorResults[0].Error {
		errInfos := make([]string, 0)
		for _, result := range errorResults {
			errInfos = append(errInfos, fmt.Sprintf(
				"[%s]: '%v' 不满足条件 '%s'",
				result.FailedField,
				result.Value,
				result.Tag,
			))
		}
		return fmt.Errorf(strings.Join(errInfos, ", "))
	}
	return nil
}
