package api

import (
	"strconv"

	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// BookkeepingBudgetApi 预算管理相关API
type BookkeepingBudgetApi struct {
	budgetService service.BookkeepingBudgetService
}

// @Summary 创建预算
// @Description 创建新的预算
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param request body dto.CreateBudgetRequest true "预算信息"
// @Success 200 {object} dto.BudgetResponse
// @Router /bk/budgets [post]
func (api *BookkeepingBudgetApi) CreateBudget(c *gin.Context) {
	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务创建预算
	result, err := api.budgetService.CreateBudget(userId, req)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取预算列表
// @Description 获取预算列表，支持分页和筛选
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param page query int true "页码" default(1)
// @Param page_size query int true "每页大小" default(10)
// @Param type query string false "预算类型 (overall, category)"
// @Param period query string false "预算周期 (weekly, monthly, yearly)"
// @Param category_id query int false "分类ID（仅当筛选分类预算时使用）"
// @Param is_active query bool false "是否激活"
// @Success 200 {object} dto.BudgetListResponse
// @Router /bk/budgets [get]
func (api *BookkeepingBudgetApi) ListBudgets(c *gin.Context) {
	var query dto.BudgetQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 默认参数设置
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取预算列表
	result, err := api.budgetService.ListBudgets(userId, query)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取预算详情
// @Description 获取单个预算的详细信息
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param id path int true "预算ID"
// @Success 200 {object} dto.BudgetResponse
// @Router /bk/budgets/{id} [get]
func (api *BookkeepingBudgetApi) GetBudget(c *gin.Context) {
	// 解析预算ID
	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorWithMsg(c, "无效的预算ID")
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取预算详情
	result, err := api.budgetService.GetBudget(userId, uint(budgetID))
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取预算进度
// @Description 获取单个预算的当前执行进度
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param id path int true "预算ID"
// @Success 200 {object} dto.BudgetProgressResponse
// @Router /bk/budgets/{id}/progress [get]
func (api *BookkeepingBudgetApi) GetBudgetProgress(c *gin.Context) {
	// 解析预算ID
	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorWithMsg(c, "无效的预算ID")
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取预算进度
	result, err := api.budgetService.GetBudgetProgress(userId, uint(budgetID))
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 更新预算
// @Description 更新预算信息
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param id path int true "预算ID"
// @Param request body dto.UpdateBudgetRequest true "更新信息"
// @Success 200 {object} dto.BudgetResponse
// @Router /bk/budgets/{id} [put]
func (api *BookkeepingBudgetApi) UpdateBudget(c *gin.Context) {
	// 解析预算ID
	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorWithMsg(c, "无效的预算ID")
		return
	}

	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务更新预算
	result, err := api.budgetService.UpdateBudget(userId, uint(budgetID), req)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 删除预算
// @Description 删除预算
// @Tags 预算管理
// @Accept json
// @Produce json
// @Param id path int true "预算ID"
// @Success 200 {object} model.common.response.Response
// @Router /bk/budgets/{id} [delete]
func (api *BookkeepingBudgetApi) DeleteBudget(c *gin.Context) {
	// 解析预算ID
	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorWithMsg(c, "无效的预算ID")
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务删除预算
	if err := api.budgetService.DeleteBudget(userId, uint(budgetID)); err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithMessage(c, "删除成功")
}

// @Summary 获取所有激活预算的进度
// @Description 获取所有激活的预算及其当前执行进度
// @Tags 预算管理
// @Accept json
// @Produce json
// @Success 200 {array} dto.BudgetProgressResponse
// @Router /bk/budgets/active-progress [get]
func (api *BookkeepingBudgetApi) ListActiveBudgetProgress(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取激活的预算进度列表
	result, err := api.budgetService.ListActiveBudgetProgress(userId)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 检查预算警告
// @Description 获取达到或超过提醒阈值的预算列表
// @Tags 预算管理
// @Accept json
// @Produce json
// @Success 200 {array} dto.BudgetProgressResponse
// @Router /bk/budgets/alerts [get]
func (api *BookkeepingBudgetApi) CheckBudgetAlerts(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取预算警告
	result, err := api.budgetService.CheckBudgetAlerts(userId)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
} 