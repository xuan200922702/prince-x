package models

type PrinceRoleDept struct {
	RoleId int `gorm:"type:int(11)"`
	DeptId int `gorm:"type:int(11)"`
}

func (PrinceRoleDept) TableName() string {
	return "prince_role_dept"
}
