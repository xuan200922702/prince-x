package system

import (
	"github.com/gin-gonic/gin"
	"prince-x/models"
	"prince-x/tools"
	"prince-x/tools/app"
)

func GetInfo(c *gin.Context) {

	var roles = make([]string, 1)
	roles[0] = tools.GetRoleName(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"

	var buttons = make([]string, 1)
	buttons[0] = "*:*:*"

	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = tools.GetRoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	if tools.GetRoleName(c) == "admin" || tools.GetRoleName(c) == "系统管理员" {
		mp["permissions"] = permissions
		mp["buttons"] = buttons
	} else {
		list, _ := RoleMenu.GetPermis()
		mp["permissions"] = list
		mp["buttons"] = list
	}

	princeuser := models.PrinceUser{}
	princeuser.UserId = tools.GetUserId(c)
	user, err := princeuser.Get()
	tools.HasError(err, "", 500)

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "http://ipc.zhangpengxuan.com/2020-07-20-%E4%BD%A0%E7%9A%84%E5%90%8D%E5%AD%971.jpeg"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["username"] = user.Username
	mp["userId"] = user.UserId
	mp["deptId"] = user.DeptId
	mp["nickname"] = user.NickName

	app.OK(c, mp, "")
}
