package api

import (
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// StatisticsAPI 统计相关API
type StatisticsAPI struct {
	statisticsService service.StatisticsService
}

// @Summary 获取收支汇总
// @Description 获取指定时间范围内的收支汇总信息
// @Tags 统计
// @Accept json
// @Produce json
// @Param request query dto.StatisticsQueryRequest true "查询参数"
// @Success 200 {object} dto.IncomeExpenseSummaryResponse
// @Router /statistics/income-expense-summary [get]
func (api *StatisticsAPI) GetIncomeExpenseSummary(c *gin.Context) {
	var req dto.StatisticsQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取收支汇总
	result, err := api.statisticsService.GetIncomeExpenseSummary(userId, req.RangeType, req.StartDate, req.EndDate)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取分类汇总
// @Description 获取指定时间范围内的分类汇总信息
// @Tags 统计
// @Accept json
// @Produce json
// @Param request query dto.CategoryStatisticsRequest true "查询参数"
// @Success 200 {array} dto.CategorySummaryItem
// @Router /statistics/category-summary [get]
func (api *StatisticsAPI) GetCategorySummary(c *gin.Context) {
	var req dto.CategoryStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 获取交易类型
	var transactionType model.TransactionType
	if req.TransactionType == "income" {
		transactionType = model.TransactionTypeIncome
	} else {
		transactionType = model.TransactionTypeExpense
	}

	// 调用服务获取分类汇总
	result, err := api.statisticsService.GetCategorySummary(userId, transactionType, req.RangeType, req.StartDate, req.EndDate)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取账户余额汇总
// @Description 获取所有账户的余额汇总信息
// @Tags 统计
// @Accept json
// @Produce json
// @Success 200 {array} dto.AccountSummaryItem
// @Router /statistics/account-summary [get]
func (api *StatisticsAPI) GetAccountSummary(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取账户汇总
	result, err := api.statisticsService.GetAccountSummary(userId)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}

// @Summary 获取月度收支趋势
// @Description 获取最近几个月的收支趋势数据
// @Tags 统计
// @Accept json
// @Produce json
// @Param months_count query int false "查询的月份数量，默认为12" default(12)
// @Success 200 {object} dto.MonthlyTrendResponse
// @Router /statistics/monthly-trend [get]
func (api *StatisticsAPI) GetMonthlyTrend(c *gin.Context) {
	// 从查询参数获取months_count，默认为"12"
	monthsCountStr := c.DefaultQuery("months_count", "12")

	// 转换为整数
	monthsCount := utils.StrToInt(monthsCountStr)

	// 如果转换失败或值不合理，使用默认值
	if monthsCount <= 0 {
		monthsCount = 12
	}

	// 获取当前用户ID
	userID, _ := c.Get("userID")
	userId := userID.(uint)

	// 调用服务获取月度趋势
	result, err := api.statisticsService.GetMonthlyTrend(userId, monthsCount)
	if err != nil {
		utils.ErrorWithMsg(c, err.Error())
		return
	}

	utils.OkWithData(c, result)
}
