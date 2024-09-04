package router

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/dotdancer/gogofly/docs"
	"github.com/dotdancer/gogofly/global"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type IFnRegisterRouter = func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup)

var (
	gfnRouters []IFnRegisterRouter
)

// RegisterRouter 注册路由回调函数
func RegisterRouter(fn IFnRegisterRouter) {
	if fn == nil {
		return
	}
	gfnRouters = append(gfnRouters, fn)
}

// InitRouter 初始化系统路由
func InitRouter() {
	//创建监听 ctrl + c,应用退出信号的上下文
	ctx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelCtx()

	//初始化gin框架，并且注册相关路由
	r := gin.Default()
	rgPublic := r.Group("/api/v1/public")
	rgAuth := r.Group("/api/v1/auth")

	//初始哈基础平台路由
	InitBasePlatformRouters()

	//自定义验证器
	//RegisterValidation()

	for _, fn := range gfnRouters {
		fn(rgPublic, rgAuth)
	}

	//集成swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	stPort := global.Config.System.Port
	if stPort == "" {
		stPort = "8099"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", stPort),
		Handler: r,
	}
	fmt.Printf("swagger address : http://localhost:%s/swagger/index.html\n", stPort)

	go func() {
		global.Logger.Infof("start http server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Logger.Error("start sever error: %s\n", err.Error())
			log.Fatalf("start sever error: %s\n", err.Error())
			return
		}
	}()
	<-ctx.Done()
	//cancelCtx()
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Error("shutdown server error: %s\n", err.Error())
		log.Fatalf("shutdown server error: %s\n", err.Error())
		return
	}

	fmt.Println("shutdown server ok")
}

func InitBasePlatformRouters() {
	InitUserRouter()
}
