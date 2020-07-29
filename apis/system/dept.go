package system

import (
	"github.com/gin-gonic/gin"
	"prince-x/models"
	"prince-x/tools"
	"prince-x/tools/app"
)

// @Summary 分页部门列表数据
// @Description 分页列表
// @Tags 部门
// @Param name query string false "name"
// @Param id query string false "id"
// @Param position query string false "position"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/deptList [get]
// @Security Bearer
func GetDeptList(c *gin.Context) {
	var Dept models.Dept
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId, _ = tools.StringToInt(c.Request.FormValue("deptId"))
	Dept.DataScope = tools.GetUserIdStr(c)
	result, err := Dept.SetDept(true)
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

func GetDeptTree(c *gin.Context) {
	var Dept models.Dept
	Dept.DeptName = c.Request.FormValue("deptName")
	Dept.Status = c.Request.FormValue("status")
	Dept.DeptId, _ = tools.StringToInt(c.Request.FormValue("deptId"))
	result, err := Dept.SetDept(false)
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}
