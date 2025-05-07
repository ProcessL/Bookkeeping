package api

import (
	"strconv"

	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// BookkeepingTransactionApi 结构体定义了交易流水管理的API处理器
type BookkeepingTransactionApi struct {
	Service service.BookkeepingTransactionService
}

// CreateTransaction godoc
// @Tags BookkeepingTransaction
// @Summary 创建交易流水
// @Description 用户创建一个新的交易记录
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   transaction_info body dto.CreateTransactionRequest true "交易信息"
// @Success 200 {object} response.Response{data=dto.TransactionResponse,msg=string} "创建成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/transactions [post]
func (a *BookkeepingTransactionApi) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c) // 从JWT获取用户ID
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	transaction, err := a.Service.CreateTransaction(userID, req)
	if err != nil {
		response.FailWithMessage("创建交易记录失败: "+err.Error(), c)
		return
	}

	response.OkWithData(transaction, c)
}

// ListTransactions godoc
// @Tags BookkeepingTransaction
// @Summary 获取交易流水列表
// @Description 获取当前用户的交易流水记录，支持分页和筛选
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   page query int false "页码，默认1"
// @Param   page_size query int false "每页数量，默认20"
// @Param   account_id query int false "账户ID筛选"
// @Param   category_id query int false "分类ID筛选"
// @Param   type query string false "交易类型筛选 (income, expense, transfer)"
// @Param   start_date query string false "开始日期筛选 (YYYY-MM-DD)"
// @Param   end_date query string false "结束日期筛选 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=dto.TransactionListResponse,msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/transactions [get]
func (a *BookkeepingTransactionApi) ListTransactions(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	accountID, _ := strconv.Atoi(c.DefaultQuery("account_id", "0"))
	categoryID, _ := strconv.Atoi(c.DefaultQuery("category_id", "0"))
	transactionType := c.DefaultQuery("type", "")
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	// 构建查询条件
	query := dto.TransactionQuery{
		Page:       page,
		PageSize:   pageSize,
		AccountID:  uint(accountID),
		CategoryID: uint(categoryID),
		Type:       transactionType,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	transactions, err := a.Service.ListTransactions(userID, query)
	if err != nil {
		response.FailWithMessage("获取交易流水列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(transactions, c)
}

// GetTransaction godoc
// @Tags BookkeepingTransaction
// @Summary 获取单个交易流水信息
// @Description 获取指定ID的交易流水详细信息
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "交易流水ID"
// @Success 200 {object} response.Response{data=dto.TransactionResponse,msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "交易记录不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/transactions/{id} [get]
func (a *BookkeepingTransactionApi) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的交易流水ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	transaction, err := a.Service.GetTransaction(userID, uint(id))
	if err != nil {
		response.FailWithMessage("获取交易流水信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(transaction, c)
}

// UpdateTransaction godoc
// @Tags BookkeepingTransaction
// @Summary 更新交易流水信息
// @Description 更新指定ID的交易流水信息
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "交易流水ID"
// @Param   transaction_info body dto.UpdateTransactionRequest true "交易信息"
// @Success 200 {object} response.Response{data=dto.TransactionResponse,msg=string} "更新成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "交易记录不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/transactions/{id} [put]
func (a *BookkeepingTransactionApi) UpdateTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的交易流水ID", c)
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	transaction, err := a.Service.UpdateTransaction(userID, uint(id), req)
	if err != nil {
		response.FailWithMessage("更新交易流水失败: "+err.Error(), c)
		return
	}

	response.OkWithData(transaction, c)
}

// DeleteTransaction godoc
// @Tags BookkeepingTransaction
// @Summary 删除交易流水
// @Description 删除指定ID的交易流水
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path int true "交易流水ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "交易记录不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/transactions/{id} [delete]
func (a *BookkeepingTransactionApi) DeleteTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.FailWithMessage("无效的交易流水ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	err = a.Service.DeleteTransaction(userID, uint(id))
	if err != nil {
		response.FailWithMessage("删除交易流水失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除交易流水成功", c)
}
