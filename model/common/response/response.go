package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	Option func(*option)
	option struct {
		code int
	}
)

// CustomerWithCode 自定义状态码
func CustomerWithCode(code int) Option {
	return func(opt *option) {
		opt.code = code
	}
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(c *gin.Context, status int, resp Response) {
	//AbortWithStatusJSON如果调用两次Result返回，会返回最后一个Result结果
	c.AbortWithStatusJSON(status, resp)
}

func Ok(c *gin.Context) {
	Result(c, http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  "操作成功",
	})
}

func OkWithMessage(c *gin.Context, msg string) {
	Result(c, http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
	})
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, Response{
		Code: SUCCESS,
		Data: data,
	})
}

func OkWithDetailed(c *gin.Context, msg string, data interface{}) {
	Result(c, http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  msg,
		Data: data,
	})
}

func Fail(c *gin.Context) {
	Result(c, http.StatusBadRequest, Response{
		Code: ERROR,
		Msg:  "操作失败",
	})
}

func FailWithMessage(c *gin.Context, msg string, opts ...Option) {
	var opt option
	for _, o := range opts {
		o(&opt)
	}
	Result(c, buildStatus(http.StatusBadRequest, opt.code), Response{
		Code: ERROR,
		Msg:  msg,
	})
}

func FailWithDetailed(c *gin.Context, msg string, data interface{}, opts ...Option) {
	var opt option
	for _, o := range opts {
		o(&opt)
	}
	Result(c, buildStatus(http.StatusBadRequest, opt.code), Response{
		Code: ERROR,
		Msg:  msg,
		Data: data,
	})
}

// buildStatus 自定义错误状态码
func buildStatus(defaultStatus int, code int) int {
	if code == 0 {
		return defaultStatus
	}
	return code
}

func AppendError(existErr, newErr error) error {
	if existErr == nil && newErr == nil {
		return nil
	}
	if existErr == nil {
		return nil
	}
	if newErr == nil {
		return nil
	}
	return fmt.Errorf("%s,%s", existErr.Error(), newErr.Error())
}
