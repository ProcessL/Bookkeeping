package utils

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); exists {
		if claimsMap, ok := claims.(map[string]interface{}); ok {
			if id, exists := claimsMap["id"]; exists {
				if idFloat, ok := id.(float64); ok {
					return uint(idFloat)
				}
			}
		}
	}
	return 0
}

// GetErrorMsg 获取验证错误的详细信息
func GetErrorMsg(obj interface{}, err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			// 获取结构体字段
			objType := reflect.TypeOf(obj)
			if objType.Kind() == reflect.Ptr {
				objType = objType.Elem()
			}
			field, _ := objType.FieldByName(e.Field())

			// 尝试获取自定义错误消息
			tagName := strings.SplitN(field.Tag.Get("binding"), ",", 2)[0]
			if tagName == "" {
				tagName = e.Tag()
			}

			// 构建错误消息
			return "字段 '" + field.Tag.Get("json") + "' " + getValidationErrorMsg(tagName, e)
		}
	}
	return err.Error()
}

// getValidationErrorMsg 根据验证标签返回对应的错误消息
func getValidationErrorMsg(tag string, e validator.FieldError) string {
	switch tag {
	case "required":
		return "不能为空"
	case "min":
		return "长度不能小于" + e.Param()
	case "max":
		return "长度不能大于" + e.Param()
	case "email":
		return "格式不正确"
	case "oneof":
		return "必须是[" + e.Param() + "]中的一个"
	default:
		return "不满足验证要求" + tag
	}
}
