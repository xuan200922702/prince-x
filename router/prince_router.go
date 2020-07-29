package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	log2 "prince-x/apis/log"
	"prince-x/apis/system"
	_ "prince-x/docs"
	"prince-x/handler"
	"prince-x/middleware"
	"prince-x/pkg/jwtauth"
	jwt "prince-x/pkg/jwtauth"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	princeBaseRouter(g)
	// 静态文件
	princeStaticFileRouter(g)

	// swagger；注意：生产环境可以注释掉
	princeSwaggerRouter(g)

	// 无需认证
	princeNoCheckRoleRouter(g)
	// 需要认证
	princeCheckRoleRouterInit(g, authMiddleware)

	return g
}

func princeBaseRouter(r *gin.RouterGroup) {
	r.GET("/", system.HelloWorld)
	r.GET("/info", handler.Ping)
}

func princeStaticFileRouter(r *gin.RouterGroup) {
	r.Static("/static", "./static")
}

func princeSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func princeNoCheckRoleRouter(r *gin.RouterGroup) {
	v1 := r.Group("/api/v1")
	v1.GET("/getCaptcha", system.GenerateCaptchaHandler)

}

func princeCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwtauth.GinJWTMiddleware) {
	r.POST("/login", authMiddleware.LoginHandler)
	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	v1 := r.Group("/api/v1")
	registerPageRouter(v1, authMiddleware)
	registerBaseRouter(v1, authMiddleware)
	registerPrinceUserRouter(v1, authMiddleware)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/getinfo", system.GetInfo)
		v1auth.POST("/logout", handler.LogOut)
	}
}

func registerPageRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{

		v1auth.GET("/princeUserList", system.GetPrinceUserList)
		v1auth.GET("/rolelist", system.GetRoleList)
		v1auth.GET("/deptList", system.GetDeptList)
		v1auth.GET("/deptTree", system.GetDeptTree)
		v1auth.GET("/rolelist", system.GetRoleList)
		v1auth.GET("/menulist", system.GetMenuList)
		v1auth.GET("/loginloglist", log2.GetLoginLogList)
	}
}

func registerPrinceUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	princeuser := v1.Group("/princeUser").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//sysuser.GET("/:userId", system.GetSysUser)
		//sysuser.GET("/", system.GetSysUserInit)
		princeuser.POST("", system.CreatePrinceUser)
		//sysuser.PUT("", system.UpdateSysUser)
		//sysuser.DELETE("/:userId", system.DeleteSysUser)
	}
}
