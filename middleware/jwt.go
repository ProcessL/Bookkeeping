package middleware

import (
	"strings"

	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.FailWithMessage(c, "无效的token")
			c.Abort()
			return
		}

		// 如果token格式为"Bearer xxx"，则提取实际的token部分
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:] // 去掉"Bearer "前缀
		}

		// 解析token
		jwt := utils.NewJWT()
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.FailWithMessage(c, "无效的token")
			c.Abort()
			return
		}

		// 将用户ID存储到上下文中
		c.Set("userID", uint(claims.ID))
		c.Next()
	}
}
