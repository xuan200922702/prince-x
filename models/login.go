package models

import (
	"prince-x/global/orm"
	"prince-x/tools"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser() (user PrinceUser, role PrinceRole, e error) {

	e = orm.Eloquent.Table("prince_user").Where("username = ? ", u.Username).Find(&user).Error
	if e != nil {
		return
	}
	_, e = tools.CompareHashAndPassword(user.Password, u.Password)
	if e != nil {
		return
	}
	e = orm.Eloquent.Table("prince_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if e != nil {
		return
	}
	return
}
