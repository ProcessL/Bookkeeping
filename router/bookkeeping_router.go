package router

import (
	"github.com/dotdancer/gogofly/api"
	"github.com/dotdancer/gogofly/middleware"
	"github.com/gin-gonic/gin"
)

// InitBookkeepingRouter 初始化记账相关模块的路由
func InitBookkeepingRouter(Router *gin.RouterGroup) {
	categoryApi := api.BookkeepingCategoryApi{}
	accountApi := api.BookkeepingAccountApi{}
	transactionApi := api.BookkeepingTransactionApi{}

	bookkeepingRouter := Router.Group("bk").Use(middleware.JWTAuth()) // 所有记账相关接口都需要认证
	{
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
	}
}
