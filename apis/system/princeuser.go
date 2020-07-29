package system

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"prince-x/models"
	"prince-x/tools"
	"prince-x/tools/app"
)

// @Summary 列表数据
// @Description 获取JSON
// @Tags 用户
// @Param username query string false "用户名"
// @Param status query string false "状态"
// @Param phone query string false "手机号"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 400 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/princeUserList [get]
// @Security Bearer
func GetPrinceUserList(c *gin.Context) {
	var data models.PrinceUser
	var err error
	var pageSize = 3
	var pageIndex = 1

	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.Username = c.Request.FormValue("userName")
	data.Status = c.Request.FormValue("status")
	data.Phone = c.Request.FormValue("phone")

	postId := c.Request.FormValue("postId")
	data.PostId, _ = tools.StringToInt(postId)

	deptId := c.Request.FormValue("deptId")
	data.DeptId, _ = tools.StringToInt(deptId)

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

// @Summary 创建用户
// @Description 获取JSON
// @Tags 用户
// @Accept  application/json
// @Product application/json
// @Param data body models.SysUser true "用户数据"
// @Success 200 {string} string	"{"code": 200, "message": "添加成功"}"
// @Success 200 {string} string	"{"code": -1, "message": "添加失败"}"
// @Router /api/v1/sysUser [post]
func CreatePrinceUser(c *gin.Context) {
	var princeuser models.PrinceUser
	err := c.BindWith(&princeuser, binding.JSON)
	tools.HasError(err, "非法数据格式", 500)

	princeuser.CreateBy = tools.GetUserIdStr(c)
	id, err := princeuser.Insert()
	tools.HasError(err, "添加失败,用户名已存在", 500)
	app.OK(c, id, "添加成功")
}
