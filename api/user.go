package api

import (
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserApi struct{}

func NewUserApi() *UserApi {
	return new(UserApi)
}

// Login
// @Tags 用户模块
// @Summary 用户登录
// @Description 用户登录详细描述
// @Accept  application/json
// @Produce  application/json
// @Param   username     body    dto.UserDTO     true   "username"
// @Success 200 {string} string	"ok"
// @Router /api/v1/public/user/login [post]
func (l *UserApi) Login(c *gin.Context) {
	var iUserDto dto.UserDTO
	err := c.ShouldBind(&iUserDto)
	if err != nil {
		if err1 := utils.ParseValidationError(err.(validator.ValidationErrors), &iUserDto); err1 != nil {
			response.FailWithMessage(c, err1.Error())
			return
		}
		response.FailWithMessage(c, err.Error())
		return
	}

	response.OkWithData(c, iUserDto)
}
