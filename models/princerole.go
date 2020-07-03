package models

import (
	"prince-x/global/orm"
	"prince-x/tools"
)

type PrinceRole struct {
	RoleId    int    `json:"roleId" gorm:"primary_key;AUTO_INCREMENT"` // 角色编码
	RoleName  string `json:"roleName" gorm:"type:varchar(128);"`       // 角色名称
	Status    string `json:"status" gorm:"type:int(1);"`               //
	RoleKey   string `json:"roleKey" gorm:"type:varchar(128);"`        //角色代码
	RoleSort  int    `json:"roleSort" gorm:"type:int(4);"`             //角色排序
	Flag      string `json:"flag" gorm:"type:varchar(128);"`           //
	CreateBy  string `json:"createBy" gorm:"type:varchar(128);"`       //
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(128);"`       //
	Remark    string `json:"remark" gorm:"type:varchar(255);"`         //备注
	Admin     bool   `json:"admin" gorm:"type:char(1);"`
	DataScope string `json:"dataScope" gorm:"type:varchar(128);"`
	Params    string `json:"params" gorm:"-"`
	MenuIds   []int  `json:"menuIds" gorm:"-"`
	DeptIds   []int  `json:"deptIds" gorm:"-"`
	BaseModel
}

func (PrinceRole) TableName() string {
	return "prince_role"
}

func (e *PrinceRole) GetPage(pageSize int, pageIndex int) ([]PrinceRole, int, error) {
	var doc []PrinceRole

	table := orm.Eloquent.Select("*").Table("prince_role")
	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}
	if e.RoleName != "" {
		table = table.Where("role_name = ?", e.RoleName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.RoleKey != "" {
		table = table.Where("role_key = ?", e.RoleKey)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("prince_role", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Order("role_sort").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

func (role *PrinceRole) Get() (SysRole PrinceRole, err error) {
	table := orm.Eloquent.Table("prince_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.First(&SysRole).Error; err != nil {
		return
	}

	return
}
