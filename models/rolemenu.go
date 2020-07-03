package models

import "prince-x/global/orm"

type RoleMenu struct {
	RoleId   int    `gorm:"type:int(11)"`
	MenuId   int    `gorm:"type:int(11)"`
	RoleName string `gorm:"type:varchar(128)"`
	CreateBy string `gorm:"type:varchar(128)"`
	UpdateBy string `gorm:"type:varchar(128)"`
}

func (RoleMenu) TableName() string {
	return "prince_role_menu"
}

func (rm *RoleMenu) GetPermis() ([]string, error) {
	var r []Menu
	table := orm.Eloquent.Select("prince_menu.permission").Table("prince_menu").Joins("left join prince_role_menu on prince_menu.menu_id = prince_role_menu.menu_id")

	table = table.Where("role_id = ?", rm.RoleId)

	table = table.Where("prince_menu.menu_type in('F','C')")
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	var list []string
	for i := 0; i < len(r); i++ {
		list = append(list, r[i].Permission)
	}
	return list, nil
}
