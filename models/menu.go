package models

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
