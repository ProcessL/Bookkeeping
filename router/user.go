package router

import (
	"github.com/dotdancer/gogofly/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitUserRouter() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		userApi := api.NewUserApi()
		rgPublicUser := rgPublic.Group("/user")
		{
			rgPublicUser.POST("/login", userApi.Login)

		}

		rgAuthUser := rgAuth.Group("/user")
		rgAuthUser.GET("", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"data": []map[string]any{
					{"id": 1, "name": "张三"},
					{"id": 2, "name": "李四"},
				},
			})
		})

		rgAuthUser.GET("/:id", func(ctx *gin.Context) {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"id":   1,
				"name": "张三",
			})
		})
	})
}
