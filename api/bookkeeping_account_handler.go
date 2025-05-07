package api

import (
	"strconv"

	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// BookkeepingAccountApi 结构体定义了账户管理的API处理器
type BookkeepingAccountApi struct {
	Service service.BookkeepingAccountService
}

// CreateAccount godoc
// @Tags BookkeepingAccount
// @Summary 创建账户
// @Description 用户创建一个新的账户
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   account_info body dto.CreateAccountRequest true "账户信息"
// @Success 200 {object} response.Response{data=dto.AccountResponse,msg=string} "创建成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/accounts [post]
func (a *BookkeepingAccountApi) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c) // 从JWT获取用户ID
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	account, err := a.Service.CreateAccount(userID, req)
	if err != nil {
		response.FailWithMessage("创建账户失败: "+err.Error(), c)
		return
	}

	response.OkWithData(account, c)
}

// ListAccounts godoc
// @Tags BookkeepingAccount
// @Summary 获取账户列表
// @Description 获取当前用户的所有账户
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Success 200 {object} response.Response{data=[]dto.AccountResponse,msg=string} "获取成功"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/accounts [get]
func (a *BookkeepingAccountApi) ListAccounts(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	accounts, err := a.Service.ListAccounts(userID)
	if err != nil {
		response.FailWithMessage("获取账户列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(accounts, c)
}

// GetAccount godoc
// @Tags BookkeepingAccount
// @Summary 获取单个账户信息
// @Description 获取指定ID的账户详细信息
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "账户ID"
// @Success 200 {object} response.Response{data=dto.AccountResponse,msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "账户不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/accounts/{id} [get]
func (a *BookkeepingAccountApi) GetAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的账户ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	account, err := a.Service.GetAccount(userID, uint(id))
	if err != nil {
		response.FailWithMessage("获取账户信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(account, c)
}

// UpdateAccount godoc
// @Tags BookkeepingAccount
// @Summary 更新账户信息
// @Description 更新指定ID的账户信息
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "账户ID"
// @Param   account_info body dto.UpdateAccountRequest true "账户信息"
// @Success 200 {object} response.Response{data=dto.AccountResponse,msg=string} "更新成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "账户不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/accounts/{id} [put]
func (a *BookkeepingAccountApi) UpdateAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的账户ID", c)
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	account, err := a.Service.UpdateAccount(userID, uint(id), req)
	if err != nil {
		response.FailWithMessage("更新账户失败: "+err.Error(), c)
		return
	}

	response.OkWithData(account, c)
}

// DeleteAccount godoc
// @Tags BookkeepingAccount
// @Summary 删除账户
// @Description 删除指定ID的账户
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "账户ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "账户不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/accounts/{id} [delete]
func (a *BookkeepingAccountApi) DeleteAccount(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的账户ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	err = a.Service.DeleteAccount(userID, uint(id))
	if err != nil {
		response.FailWithMessage("删除账户失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除账户成功", c)
}
