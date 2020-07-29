package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"prince-x/global/orm"
	"prince-x/tools"
)

type UserName struct {
	Username string `gorm:"type:varchar(64)" json:"username"`
}

type PassWord struct {
	// 密码
	Password string `gorm:"type:varchar(128)" json:"password"`
}

type LoginM struct {
	UserName
	PassWord
}
type PrinceUserView struct {
	PrinceUserId
	PrinceUserB
	LoginM
	RoleName string `gorm:"column:role_name"  json:"role_name"`
}

type PrinceUserId struct {
	UserId int `gorm:"primary_key;AUTO_INCREMENT"  json:"userId"` // 编码
}

type PrinceUserB struct {
	NickName  string `gorm:"type:varchar(128)" json:"nickName"` // 昵称
	Phone     string `gorm:"type:varchar(11)" json:"phone"`     // 手机号
	RoleId    int    `gorm:"type:int(11)" json:"roleId"`        // 角色编码
	Salt      string `gorm:"type:varchar(255)" json:"salt"`     //盐
	Avatar    string `gorm:"type:varchar(255)" json:"avatar"`   //头像
	Sex       string `gorm:"type:varchar(255)" json:"sex"`      //性别
	Email     string `gorm:"type:varchar(128)" json:"email"`    //邮箱
	DeptId    int    `gorm:"type:int(11)" json:"deptId"`        //部门编码
	PostId    int    `gorm:"type:int(11)" json:"postId"`        //职位编码
	CreateBy  string `gorm:"type:varchar(128)" json:"createBy"` //
	UpdateBy  string `gorm:"type:varchar(128)" json:"updateBy"` //
	Remark    string `gorm:"type:varchar(255)" json:"remark"`   //备注
	Status    string `gorm:"type:int(1);" json:"status"`
	DataScope string `gorm:"-" json:"dataScope"`
	Params    string `gorm:"-" json:"params"`

	BaseModel
}

type PrinceUser struct {
	PrinceUserId
	PrinceUserB
	LoginM
}

func (PrinceUser) TableName() string {
	return "prince_user"
}

type PrinceUserPage struct {
	PrinceUserId
	PrinceUserB
	LoginM
	DeptName string `gorm:"-" json:"deptName"`
}

func (e *PrinceUser) GetPage(pageSize int, pageIndex int) ([]PrinceUserPage, int, error) {
	var doc []PrinceUserPage
	table := orm.Eloquent.Select("prince_user.*,prince_dept.dept_name").Table(e.TableName())
	table = table.Joins("left join prince_dept on prince_dept.dept_id = prince_user.dept_id")

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}
	if e.Status != "" {
		table = table.Where("prince_user.status = ?", e.Status)
	}

	if e.Phone != "" {
		table = table.Where("prince_user.phone = ?", e.Phone)
	}

	if e.DeptId != 0 {
		table = table.Where("prince_user.dept_id in (select dept_id from prince_dept where dept_path like ? )", "%"+tools.IntToString(e.DeptId)+"%")
	}

	// 数据权限控制
	dataPermission := new(DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("prince_user", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("prince_user.deleted_at IS NULL").Count(&count)
	return doc, count, nil
}

// 获取用户数据
func (e *PrinceUser) Get() (PrinceUserView PrinceUserView, err error) {

	table := orm.Eloquent.Table(e.TableName()).Select([]string{"prince_user.*", "prince_role.role_name"})
	table = table.Joins("left join prince_role on prince_user.role_id=prince_role.role_id")
	if e.UserId != 0 {
		table = table.Where("user_id = ?", e.UserId)
	}

	if e.Username != "" {
		table = table.Where("username = ?", e.Username)
	}

	if e.Password != "" {
		table = table.Where("password = ?", e.Password)
	}

	if e.RoleId != 0 {
		table = table.Where("role_id = ?", e.RoleId)
	}

	if e.DeptId != 0 {
		table = table.Where("dept_id = ?", e.DeptId)
	}

	if e.PostId != 0 {
		table = table.Where("post_id = ?", e.PostId)
	}

	if err = table.First(&PrinceUserView).Error; err != nil {
		return
	}
	PrinceUserView.Password = ""
	return
}

//加密
func (e *PrinceUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

//添加
func (e PrinceUser) Insert() (id int, err error) {
	if err = e.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int
	orm.Eloquent.Table(e.TableName()).Where("username = ?", e.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = orm.Eloquent.Table(e.TableName()).Create(&e).Error; err != nil {
		return
	}
	id = e.UserId
	return
}
