package main

import (
	"github.com/dotdancer/gogofly/cmd"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

//go:generate swag init
//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy

// @title                       GoGoFly Swagger API接口文档
// @version                     v1.0.1
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	defer cmd.Clear()
	cmd.Start()
	//r := gin.New()
	//r.Use(Logger())
	//
	//r.GET("/test", func(c *gin.Context) {
	//	example := c.MustGet("example").(string)
	//
	//	// it would print: "12345"
	//	log.Println("example:", example)
	//})
	//
	//// Listen and serve on 0.0.0.0:8080
	//r.Run(":8000")
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Println("latency:", latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println("status:", status)
	}
}
