package system

import (
	"github.com/gin-gonic/gin"
	"prince-x/models"
	"prince-x/tools"
	"prince-x/tools/app"
)

// @Summary Menu列表数据
// @Description 获取JSON
// @Tags 菜单
// @Param menuName query string false "menuName"
// @Success 200 {string} string "{"code": 200, "data": [...]}"
// @Success 200 {string} string "{"code": -1, "message": "抱歉未找到相关信息"}"
// @Router /api/v1/menulist [get]
// @Security Bearer
func GetMenuList(c *gin.Context) {
	var Menu models.Menu
	Menu.MenuName = c.Request.FormValue("menuName")
	Menu.Visible = c.Request.FormValue("visible")
	Menu.Title = c.Request.FormValue("title")
	Menu.DataScope = tools.GetUserIdStr(c)
	result, err := Menu.SetMenu()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}
