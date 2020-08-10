package models

import (
	"fmt"
	"prince-x/global/orm"
	"prince-x/tools"
)

type RoleMenu struct {
	RoleId   int    `gorm:"type:int(11)" json:"role_id"`
	MenuId   int    `gorm:"type:int(11)" json:"menu_id"`
	RoleName string `gorm:"type:varchar(128)" json:"role_name"`
	CreateBy string `gorm:"type:varchar(128)" json:"create_by"`
	UpdateBy string `gorm:"type:varchar(128)" json:"update_by"`
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

type MenuPath struct {
	Path string `json:"path"`
}

func (rm *RoleMenu) Get() ([]RoleMenu, error) {
	var r []RoleMenu
	table := orm.Eloquent.Table("prince_role_menu")
	if rm.RoleId != 0 {
		table = table.Where("role_id = ?", rm.RoleId)

	}
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) GetIDS() ([]MenuPath, error) {
	var r []MenuPath
	table := orm.Eloquent.Select("prince_menu.path").Table("prince_role_menu")
	table = table.Joins("left join prince_role on prince_role.role_id=prince_role_menu.role_id")
	table = table.Joins("left join prince_menu on prince_menu.id=prince_role_menu.menu_id")
	table = table.Where("prince_role.role_name = ? and prince_menu.type=1", rm.RoleName)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rm *RoleMenu) DeleteRoleMenu(roleId int) (bool, error) {
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := tx.Table("prince_role_dept").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Table("prince_role_menu").Where("role_id = ?", roleId).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	var role PrinceRole
	if err := tx.Table("prince_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	sql3 := "delete from casbin_rule where v0= '" + role.RoleKey + "';"
	if err := tx.Exec(sql3).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil

}

func (rm *RoleMenu) BatchDeleteRoleMenu(roleIds []int) (bool, error) {
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}

	if err := tx.Table("prince_role_menu").Where("role_id in (?)", roleIds).Delete(&rm).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	var role []PrinceRole
	if err := tx.Table("prince_role").Where("role_id in (?)", roleIds).Find(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	sql := ""
	for i := 0; i < len(role); i++ {
		sql += "delete from casbin_rule where v0= '" + role[i].RoleKey + "';"
	}
	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil

}

func (rm *RoleMenu) Insert(roleId int, menuId []int) (bool, error) {
	var role PrinceRole
	tx := orm.Eloquent.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false, err
	}
	if err := tx.Table("prince_role").Where("role_id = ?", roleId).First(&role).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	var menu []Menu
	if err := tx.Table("prince_menu").Where("menu_id in (?)", menuId).Find(&menu).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	//ORM不支持批量插入所以需要拼接 sql 串
	sql := "INSERT INTO `prince_role_menu` (`role_id`,`menu_id`,`role_name`) VALUES "

	sql2 := "INSERT INTO `casbin_rule`  (`p_type`,`v0`,`v1`,`v2`) VALUES "
	for i := 0; i < len(menu); i++ {
		if len(menu)-1 == i {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,%d,'%s');", role.RoleId, menu[i].MenuId, role.RoleKey)
			if menu[i].MenuType == "A" {
				sql2 += fmt.Sprintf("('p','%s','%s','%s');", role.RoleKey, menu[i].Path, menu[i].Action)
			}
		} else {
			sql += fmt.Sprintf("(%d,%d,'%s'),", role.RoleId, menu[i].MenuId, role.RoleKey)
			if menu[i].MenuType == "A" {
				sql2 += fmt.Sprintf("('p','%s','%s','%s'),", role.RoleKey, menu[i].Path, menu[i].Action)
			}
		}
	}
	if err := tx.Exec(sql).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	sql2 = sql2[0:len(sql2)-1] + ";"
	if err := tx.Exec(sql2).Error; err != nil {
		tx.Rollback()
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil
}

func (rm *RoleMenu) Delete(RoleId string, MenuID string) (bool, error) {
	rm.RoleId, _ = tools.StringToInt(RoleId)
	table := orm.Eloquent.Table("prince_role_menu").Where("role_id = ?", RoleId)
	if MenuID != "" {
		table = table.Where("menu_id = ?", MenuID)
	}
	if err := table.Delete(&rm).Error; err != nil {
		return false, err
	}
	return true, nil

}
