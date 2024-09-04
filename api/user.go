package api

import (
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserInfoApi struct{}

var userInfoService = new(service.UserInfoService)

// Login
// @Tags 用户模块
// @Summary 用户登录
// @Description 用户登录详细描述
// @Accept  application/json
// @Produce  application/json
// @Param   username     body    dto.UserInfo     true   "username"
// @Success 200 {string} string	"ok"
// @Router /api/v1/public/user/login [post]
func (l *UserInfoApi) Login(c *gin.Context) {
	var iUserDto dto.UserInfo
	err := c.ShouldBind(&iUserDto)
	if err != nil {
		if errs := utils.ParseValidationError(err.(validator.ValidationErrors), &iUserDto); errs != nil {
			response.FailWithMessage(c, errs.Error(), response.CustomerWithCode(http.StatusBadRequest))
			return
		}
		response.FailWithMessage(c, err.Error())
		return
	}
	err = userInfoService.AddUser(iUserDto)
	if err != nil {
		response.FailWithMessage(c, err.Error())
	}
	response.Ok(c)
}
