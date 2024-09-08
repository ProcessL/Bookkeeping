package api

import (
	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type UserInfoApi struct{}

var userInfoService = new(service.UserInfoService)

// Login
// @Tags 用户模块
// @Summary 添加用户
// @Description 用户注册详细描述
// @Accept  application/json
// @Produce  application/json
// @Param   username     formDate    dto.LoginUserDto     true   "username"
// @Success 200 {string} string	"ok"
// @Router /api/v1/public/user/login [post]
func (l *UserInfoApi) Login(c *gin.Context) {
	var iUserDto dto.LoginUserDto
	err := c.ShouldBind(&iUserDto)
	if err != nil {
		if errs := utils.ParseValidationError(err.(validator.ValidationErrors), &iUserDto); errs != nil {
			response.FailWithMessage(c, errs.Error(), response.CustomerWithCode(http.StatusBadRequest))
			return
		}
		response.FailWithMessage(c, err.Error())
		return
	}
	user, token, err := userInfoService.Login(&iUserDto)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	if global.Redis.Set(c, strconv.Itoa(int(user.ID))+user.Username, token, 60*time.Second).Err() != nil {
		response.FailWithDetailed(c, "redis set token error", nil, response.CustomerWithCode(http.StatusInternalServerError))
		return
	}
	response.OkWithData(c, token)
}

// AddUser
// @Tags 用户模块
// @Summary 添加用户
// @Description 用户注册详细描述
// @Accept  application/json
// @Produce  application/json
// @Param   username     body    dto.AddUserDto     true   "username"
// @Success 200 {string} string	"ok"
// @Router /api/v1/auth/user/addUser [post]
func (l *UserInfoApi) AddUser(c *gin.Context) {
	var iUserDto dto.AddUserDto
	err := c.ShouldBind(&iUserDto)
	if err != nil {
		if errs := utils.ParseValidationError(err.(validator.ValidationErrors), &iUserDto); errs != nil {
			response.FailWithMessage(c, errs.Error(), response.CustomerWithCode(http.StatusBadRequest))
			return
		}
		response.FailWithMessage(c, err.Error())
		return
	}
	err = userInfoService.AddUser(&iUserDto)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.Ok(c)
}

// GetUserById
// @Tags 用户模块
// @Summary 获取用户详情
// @Description 获取用户详情描述
// @Accept  application/json
// @Produce  application/json
// @Param   id     path    int     true   "username"
// @Success 200 {string} string	"ok"
// @Router /api/v1/auth/user/getUserById/{id} [get]
func (l *UserInfoApi) GetUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := userInfoService.GetUserById(id)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}
	response.OkWithData(c, user)
}
