package models

import (
	"errors"
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

type MenuIdList struct {
	MenuId int `json:"menuId"`
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

func (role *PrinceRole) Get() (PrinceRole PrinceRole, err error) {
	table := orm.Eloquent.Table("prince_role")
	if role.RoleId != 0 {
		table = table.Where("role_id = ?", role.RoleId)
	}
	if role.RoleName != "" {
		table = table.Where("role_name = ?", role.RoleName)
	}
	if err = table.First(&PrinceRole).Error; err != nil {
		return
	}

	return
}

type DeptIdList struct {
	DeptId int `json:"DeptId"`
}

func (role *PrinceRole) GetRoleDeptId() ([]int, error) {
	deptIds := make([]int, 0)
	deptList := make([]DeptIdList, 0)
	if err := orm.Eloquent.Table("prince_role_dept").Select("prince_role_dept.dept_id").Joins("LEFT JOIN prince_dept on prince_dept.dept_id=prince_role_dept.dept_id").Where("role_id = ? ", role.RoleId).Where(" prince_role_dept.dept_id not in(select prince_dept.parent_id from prince_role_dept LEFT JOIN prince_dept on prince_dept.dept_id=prince_role_dept.dept_id where role_id =? )", role.RoleId).Find(&deptList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(deptList); i++ {
		deptIds = append(deptIds, deptList[i].DeptId)
	}

	return deptIds, nil
}

// 获取角色对应的菜单ids
func (role *PrinceRole) GetRoleMeunId() ([]int, error) {
	menuIds := make([]int, 0)
	menuList := make([]MenuIdList, 0)
	if err := orm.Eloquent.Table("prince_role_menu").Select("prince_role_menu.menu_id").Where("role_id = ? ", role.RoleId).Where(" prince_role_menu.menu_id not in(select prince_menu.parent_id from prince_role_menu LEFT JOIN prince_menu on prince_menu.menu_id=prince_role_menu.menu_id where role_id =?  and parent_id is not null)", role.RoleId).Find(&menuList).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (role *PrinceRole) Insert() (id int, err error) {
	i := 0
	orm.Eloquent.Table(role.TableName()).Where("role_name=? or role_key = ?", role.RoleName, role.RoleKey).Count(&i)
	if i > 0 {
		return 0, errors.New("角色名称或者角色标识已经存在！")
	}
	role.UpdateBy = ""
	result := orm.Eloquent.Table(role.TableName()).Create(&role)
	if result.Error != nil {
		err = result.Error
		return
	}
	id = role.RoleId
	return
}

//修改
func (role *PrinceRole) Update(id int) (update PrinceRole, err error) {
	if err = orm.Eloquent.Table(role.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if role.RoleName != "" && role.RoleName != update.RoleName {
		return update, errors.New("角色名称不允许修改！")
	}

	if role.RoleKey != "" && role.RoleKey != update.RoleKey {
		return update, errors.New("角色标识不允许修改！")
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(role.TableName()).Model(&update).Updates(&role).Error; err != nil {
		return
	}
	return
}

//批量删除
func (role *PrinceRole) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Unscoped().Where("role_id in (?)", id).Delete(PrinceRole{}).Error; err != nil {
		return
	}
	Result = true
	return
}
