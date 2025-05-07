package utils

import (
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/gin-gonic/gin"
)

// HandleValidationError 处理请求校验错误
func HandleValidationError(c *gin.Context, err error) {
	// 从请求中获取需要验证的数据类型，但不处理错误
	// GetErrorMsg 需要一个具体对象和错误信息
	ErrorWithMsg(c, err.Error())
}

// ErrorWithMsg 返回错误消息
func ErrorWithMsg(c *gin.Context, msg string) {
	response.FailWithMessage(c, msg)
}

// OkWithData 返回成功数据
func OkWithData(c *gin.Context, data interface{}) {
	response.OkWithData(c, data)
}

// OkWithMessage 返回成功消息
func OkWithMessage(c *gin.Context, msg string) {
	response.OkWithMessage(c, msg)
}

// OkWithDetailed 返回成功消息和数据
func OkWithDetailed(c *gin.Context, msg string, data interface{}) {
	response.OkWithDetailed(c, msg, data)
}
