package models

import (
	"prince-x/global/orm"
	"prince-x/tools"
)

type Menu struct {
	MenuId     int    `json:"menuId" gorm:"primary_key;AUTO_INCREMENT"`
	MenuName   string `json:"menuName" gorm:"type:varchar(128);"`
	Title      string `json:"title" gorm:"type:varchar(64);"`
	Icon       string `json:"icon" gorm:"type:varchar(128);"`
	Path       string `json:"path" gorm:"type:varchar(128);"`
	Paths      string `json:"paths" gorm:"type:varchar(128);"`
	MenuType   string `json:"menuType" gorm:"type:varchar(1);"`
	Action     string `json:"action" gorm:"type:varchar(16);"`
	Permission string `json:"permission" gorm:"type:varchar(32);"`
	ParentId   int    `json:"parentId" gorm:"type:int(11);"`
	NoCache    bool   `json:"noCache" gorm:"type:char(1);"`
	Breadcrumb string `json:"breadcrumb" gorm:"type:varchar(255);"`
	Component  string `json:"component" gorm:"type:varchar(255);"`
	Sort       int    `json:"sort" gorm:"type:int(4);"`
	Visible    string `json:"visible" gorm:"type:char(1);"`
	CreateBy   string `json:"createBy" gorm:"type:varchar(128);"`
	UpdateBy   string `json:"updateBy" gorm:"type:varchar(128);"`
	IsFrame    string `json:"isFrame" gorm:"type:int(1);DEFAULT:0;"`
	DataScope  string `json:"dataScope" gorm:"-"`
	Params     string `json:"params" gorm:"-"`
	RoleId     int    `gorm:"-"`
	Children   []Menu `json:"children" gorm:"-"`
	IsSelect   bool   `json:"is_select" gorm:"-"`
	BaseModel
}

func (Menu) TableName() string {
	return "prince_menu"
}

func (e *Menu) SetMenu() (m []Menu, err error) {
	menulist, err := e.GetPage()

	m = make([]Menu, 0)
	for i := 0; i < len(menulist); i++ {
		if menulist[i].ParentId != 0 {
			continue
		}
		menusInfo := DiguiMenu(&menulist, menulist[i])

		m = append(m, menusInfo)
	}
	return
}

func (e *Menu) GetPage() (Menus []Menu, err error) {
	table := orm.Eloquent.Table(e.TableName())
	if e.MenuName != "" {
		table = table.Where("menu_name = ?", e.MenuName)
	}
	if e.Title != "" {
		table = table.Where("title = ?", e.Title)
	}
	if e.Visible != "" {
		table = table.Where("visible = ?", e.Visible)
	}
	if e.MenuType != "" {
		table = table.Where("menu_type = ?", e.MenuType)
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err = dataPermission.GetDataScope("prince_menu", table)
	if err != nil {
		return nil, err
	}
	if err = table.Order("sort").Find(&Menus).Error; err != nil {
		return
	}
	return
}

func DiguiMenu(menulist *[]Menu, menu Menu) Menu {
	list := *menulist

	min := make([]Menu, 0)
	for j := 0; j < len(list); j++ {

		if menu.MenuId != list[j].ParentId {
			continue
		}
		mi := Menu{}
		mi.MenuId = list[j].MenuId
		mi.MenuName = list[j].MenuName
		mi.Title = list[j].Title
		mi.Icon = list[j].Icon
		mi.Path = list[j].Path
		mi.MenuType = list[j].MenuType
		mi.Action = list[j].Action
		mi.Permission = list[j].Permission
		mi.ParentId = list[j].ParentId
		mi.NoCache = list[j].NoCache
		mi.Breadcrumb = list[j].Breadcrumb
		mi.Component = list[j].Component
		mi.Sort = list[j].Sort
		mi.Visible = list[j].Visible
		mi.CreatedAt = list[j].CreatedAt
		mi.Children = []Menu{}

		if mi.MenuType != "F" {
			ms := DiguiMenu(menulist, mi)
			min = append(min, ms)

		} else {
			min = append(min, mi)
		}

	}
	menu.Children = min
	return menu
}
