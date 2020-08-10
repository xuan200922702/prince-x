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
	registerUserCenterRouter(v1, authMiddleware)
	registerRoleRouter(v1, authMiddleware)
	registerDeptRouter(v1, authMiddleware)
}

func registerUserCenterRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	user := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//user.GET("/profile", system.GetSysUserProfile)
		//user.POST("/avatar", system.InsetSysUserAvatar)
		user.PUT("/pwd", system.PrinceUserUpdatePwd)
	}
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/getinfo", system.GetInfo)
		v1auth.GET("/menuids", system.GetMenuIDS)
		v1auth.POST("/logout", handler.LogOut)
	}
}

//获取列表
func registerPageRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	v1auth := v1.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{

		v1auth.GET("/princeUserList", system.GetPrinceUserList)
		v1auth.GET("/rolelist", system.GetRoleList)
		v1auth.GET("/deptList", system.GetDeptList)
		v1auth.GET("/deptTree", system.GetDeptTree)
		v1auth.GET("/menulist", system.GetMenuList)
		v1auth.GET("/loginloglist", log2.GetLoginLogList)
	}
}

//princeUser
func registerPrinceUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	princeuser := v1.Group("/princeUser").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		//sysuser.GET("/:userId", system.GetSysUser)
		//sysuser.GET("/", system.GetSysUserInit)
		princeuser.POST("", system.CreatePrinceUser)
		princeuser.PUT("", system.UpdatePrinceUser)
		princeuser.DELETE("/:userId", system.DeletePrinceUser)
	}
}

//role
func registerRoleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	role := v1.Group("/princeRole").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		role.GET("/:roleId", system.GetRole)
		role.POST("", system.InsertRole)
		role.PUT("", system.UpdateRole)
		role.DELETE("/:roleId", system.DeleteRole)
	}
}

//dept
func registerDeptRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	dept := v1.Group("/dept").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		dept.GET("/:deptId", system.GetDept)
		dept.POST("", system.InsertDept)
		dept.PUT("", system.UpdateDept)
		dept.DELETE("/:id", system.DeleteDept)
	}
}
