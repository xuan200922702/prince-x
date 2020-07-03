package system

import (
	"github.com/gin-gonic/gin"
	"prince-x/models"
	"prince-x/tools"
	"prince-x/tools/app"
)

// @Summary 角色列表数据
// @Description Get JSON
// @Tags 角色/Role
// @Param roleName query string false "roleName"
// @Param status query string false "status"
// @Param roleKey query string false "roleKey"
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/rolelist [get]
// @Security
func GetRoleList(c *gin.Context) {
	var data models.PrinceRole
	var err error
	var pageSize = 10
	var pageIndex = 1

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}
