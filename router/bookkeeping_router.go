package router

import (
	"github.com/dotdancer/gogofly/api"
	"github.com/dotdancer/gogofly/middleware"
	"github.com/gin-gonic/gin"
)

// InitBookkeepingRouter 初始化记账相关模块的路由
func InitBookkeepingRouter() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		categoryApi := api.BookkeepingCategoryApi{}
		accountApi := api.BookkeepingAccountApi{}
		transactionApi := api.BookkeepingTransactionApi{}
		statisticsApi := api.StatisticsAPI{}
		budgetApi := api.BookkeepingBudgetApi{}

		// 所有记账相关接口都需要认证
		bookkeepingRouter := rgAuth.Group("bk")
		bookkeepingRouter.Use(middleware.JWTAuth())

		// 分类管理路由
		categoryRouter := bookkeepingRouter.Group("categories")
		{
			categoryRouter.POST("", categoryApi.CreateCategory)            // 创建分类
			categoryRouter.GET("", categoryApi.ListCategories)             // 获取分类列表 (层级)
			categoryRouter.GET("/flat", categoryApi.ListAllCategoriesFlat) // 获取所有分类列表 (扁平)
			categoryRouter.GET("/:id", categoryApi.GetCategory)            // 获取单个分类信息
			categoryRouter.PUT("/:id", categoryApi.UpdateCategory)         // 更新分类信息
			categoryRouter.DELETE("/:id", categoryApi.DeleteCategory)      // 删除分类
		}

		// 账户管理路由
		accountRouter := bookkeepingRouter.Group("accounts")
		{
			accountRouter.POST("", accountApi.CreateAccount)       // 创建账户
			accountRouter.GET("", accountApi.ListAccounts)         // 获取账户列表
			accountRouter.GET("/:id", accountApi.GetAccount)       // 获取单个账户信息
			accountRouter.PUT("/:id", accountApi.UpdateAccount)    // 更新账户信息
			accountRouter.DELETE("/:id", accountApi.DeleteAccount) // 删除账户
		}

		// 交易流水管理路由
		transactionRouter := bookkeepingRouter.Group("transactions")
		{
			transactionRouter.POST("", transactionApi.CreateTransaction)       // 创建交易流水
			transactionRouter.GET("", transactionApi.ListTransactions)         // 获取交易流水列表
			transactionRouter.GET("/:id", transactionApi.GetTransaction)       // 获取单个交易流水信息
			transactionRouter.PUT("/:id", transactionApi.UpdateTransaction)    // 更新交易流水信息
			transactionRouter.DELETE("/:id", transactionApi.DeleteTransaction) // 删除交易流水
		}

		// 统计分析路由
		statisticsRouter := bookkeepingRouter.Group("statistics")
		{
			statisticsRouter.GET("/income-expense-summary", statisticsApi.GetIncomeExpenseSummary) // 收支汇总
			statisticsRouter.GET("/category-summary", statisticsApi.GetCategorySummary)            // 分类汇总
			statisticsRouter.GET("/account-summary", statisticsApi.GetAccountSummary)              // 账户余额汇总
			statisticsRouter.GET("/monthly-trend", statisticsApi.GetMonthlyTrend)                  // 月度收支趋势
		}

		// 预算管理路由
		budgetRouter := bookkeepingRouter.Group("budgets")
		{
			budgetRouter.POST("", budgetApi.CreateBudget)                            // 创建预算
			budgetRouter.GET("", budgetApi.ListBudgets)                              // 获取预算列表
			budgetRouter.GET("/active-progress", budgetApi.ListActiveBudgetProgress) // 获取激活的预算进度列表
			budgetRouter.GET("/alerts", budgetApi.CheckBudgetAlerts)                 // 获取预算提醒
			budgetRouter.GET("/:id", budgetApi.GetBudget)                            // 获取单个预算信息
			budgetRouter.GET("/:id/progress", budgetApi.GetBudgetProgress)           // 获取预算进度
			budgetRouter.PUT("/:id", budgetApi.UpdateBudget)                         // 更新预算信息
			budgetRouter.DELETE("/:id", budgetApi.DeleteBudget)                      // 删除预算
		}
	})
}
