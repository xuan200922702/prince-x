package router

import (
	"github.com/gin-gonic/gin"
	"prince-x/handler"
	"prince-x/middleware"
	_ "prince-x/pkg/jwtauth"
	"prince-x/tools"
	config2 "prince-x/tools/config"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	if config2.ApplicationConfig.IsHttps {
		r.Use(handler.TlsHandler())
	}
	middleware.InitMiddleware(r)
	// TODO: the jwt middleware
	authMiddleware, err := middleware.AuthInit()
	tools.HasError(err, "JWT Init Error", 500)

	// TODO: 注册系统路由
	InitSysRouter(r, authMiddleware)

	// TODO: 注册业务路由

	//TODO: 这里可存放业务路由，里边并无实际路由是有演示代码
	//InitExamplesRouter(r, authMiddleware)

	return r
}
