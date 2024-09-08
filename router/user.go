package router

import (
	"github.com/dotdancer/gogofly/api"
	"github.com/gin-gonic/gin"
)

func InitUserRouter() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		userApi := api.UserInfoApi{}
		rgPublicUser := rgPublic.Group("/user")
		{
			rgPublicUser.POST("/login", userApi.Login)
		}

		rgAuthUser := rgAuth.Group("/user")
		{
			rgAuthUser.POST("/addUser", userApi.AddUser)
			rgAuthUser.GET("/getUserById/:id", userApi.GetUserById)
		}
	})
}
