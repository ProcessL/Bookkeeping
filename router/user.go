package router

import (
	"github.com/dotdancer/gogofly/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitUserRouter() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		userApi := api.UserInfoApi{}
		rgPublicUser := rgPublic.Group("/user")
		//rgPublicUser.Use(func(c *gin.Context) {
		//	header := c.GetHeader("x-apikey")
		//	fmt.Println("header:", header)
		//	c.Next()
		//})
		{
			rgPublicUser.POST("/login", userApi.Login)
			rgPublicUser.GET("/dome", func(c *gin.Context) {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"msg": "登录页面", "url": "/login.html",
					"data": gin.H{
						"username": "admin", "password": "123456",
					},
				})
			})
			rgPublicUser.GET("/createTenableData", userApi.CreateTenableData)
			rgPublicUser.GET("/scanResult", userApi.ScanResult)
			rgPublicUser.GET("/analysis/:id", userApi.Analysis)
		}

		rgAuthUser := rgAuth.Group("/user")
		{
			rgAuthUser.POST("/addUser", userApi.AddUser)
			rgAuthUser.GET("/getUserById/:id", userApi.GetUserById)
		}
	})
}
