package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"prince-x/tools"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

func (e *DataPermission) GetDataScope(tbname string, table *gorm.DB) (*gorm.DB, error) {
	PrinceUser := new(PrinceUser)
	PrinceRole := new(PrinceRole)
	PrinceUser.UserId = e.UserId
	user, err := PrinceUser.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	PrinceRole.RoleId = user.RoleId
	role, err := PrinceRole.Get()
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
		table = table.Where(tbname+".create_by in (select prince_user.user_id from prince_role_dept left join prince_user on prince_user.dept_id=prince_role_dept.dept_id where prince_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		table = table.Where(tbname+".create_by in (SELECT user_id from prince_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		table = table.Where(tbname+".create_by in (SELECT user_id from prince_user where prince_user.dept_id in(select dept_id from prince_dept where dept_path like ? ))", "%"+tools.IntToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		table = table.Where(tbname+".create_by = ?", e.UserId)
	}

	return table, nil
}
