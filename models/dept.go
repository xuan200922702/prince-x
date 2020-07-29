package models

import (
	"prince-x/global/orm"
	"prince-x/tools"
)

type Dept struct {
	DeptId    int    `json:"deptId" gorm:"primary_key;AUTO_INCREMENT"` //部门编码
	ParentId  int    `json:"parentId" gorm:"type:int(11);"`            //上级部门
	DeptPath  string `json:"deptPath" gorm:"type:varchar(255);"`       //
	DeptName  string `json:"deptName"  gorm:"type:varchar(128);"`      //部门名称
	Sort      int    `json:"sort" gorm:"type:int(4);"`                 //排序
	Leader    string `json:"leader" gorm:"type:varchar(128);"`         //负责人
	Phone     string `json:"phone" gorm:"type:varchar(11);"`           //手机
	Email     string `json:"email" gorm:"type:varchar(64);"`           //邮箱
	Status    string `json:"status" gorm:"type:int(1);"`               //状态
	CreateBy  string `json:"createBy" gorm:"type:varchar(64);"`
	UpdateBy  string `json:"updateBy" gorm:"type:varchar(64);"`
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params" gorm:"-"`
	Children  []Dept `json:"children" gorm:"-"`
	BaseModel
}

func (Dept) TableName() string {
	return "prince_dept"
}

func (e *Dept) SetDept(bl bool) ([]Dept, error) {
	list, err := e.GetPage(bl)

	m := make([]Dept, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		info := Digui(&list, list[i])

		m = append(m, info)
	}
	return m, err
}

func (e *Dept) GetPage(bl bool) ([]Dept, error) {
	var doc []Dept

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}
	if e.DeptName != "" {
		table = table.Where("dept_name = ?", e.DeptName)
	}
	if e.Status != "" {
		table = table.Where("status = ?", e.Status)
	}
	if e.DeptPath != "" {
		table = table.Where("deptPath like %?%", e.DeptPath)
	}
	if bl {
		// 数据权限控制
		dataPermission := new(DataPermission)
		dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
		tableper, err := dataPermission.GetDataScope("prince_dept", table)
		if err != nil {
			return nil, err
		}
		table = tableper
	}

	if err := table.Order("sort").Find(&doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func Digui(deptlist *[]Dept, menu Dept) Dept {
	list := *deptlist

	min := make([]Dept, 0)
	for j := 0; j < len(list); j++ {

		if menu.DeptId != list[j].ParentId {
			continue
		}
		mi := Dept{}
		mi.DeptId = list[j].DeptId
		mi.ParentId = list[j].ParentId
		mi.DeptPath = list[j].DeptPath
		mi.DeptName = list[j].DeptName
		mi.Sort = list[j].Sort
		mi.Leader = list[j].Leader
		mi.Phone = list[j].Phone
		mi.Email = list[j].Email
		mi.Status = list[j].Status
		mi.Children = []Dept{}
		ms := Digui(deptlist, mi)
		min = append(min, ms)

	}
	menu.Children = min
	return menu
}
