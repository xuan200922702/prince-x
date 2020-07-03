package main

import "prince-x/cmd"

// @title prince-x API
// @version 0.0.1
// @description 基于Gin + Vue + Element UI的前后端分离权限管理系统的接口文档
// @description 添加qq: 200922702 请备注，谢谢！
// @license.name Zhang Pengxuan
// @license.url https://www.baidu.com

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()

}
