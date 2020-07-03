package models

import (
	"prince-x/global/orm"
	"time"
)

type LoginLog struct {
	InfoId        int       `json:"infoId" gorm:"primary_key;AUTO_INCREMENT"` //主键
	Username      string    `json:"username" gorm:"type:varchar(128);"`       //用户名
	Status        string    `json:"status" gorm:"type:int(1);"`               //状态
	Ipaddr        string    `json:"ipaddr" gorm:"type:varchar(255);"`         //ip地址
	LoginLocation string    `json:"loginLocation" gorm:"type:varchar(255);"`  //归属地
	Browser       string    `json:"browser" gorm:"type:varchar(255);"`        //浏览器
	Os            string    `json:"os" gorm:"type:varchar(255);"`             //系统
	Platform      string    `json:"platform" gorm:"type:varchar(255);"`       // 固件
	LoginTime     time.Time `json:"loginTime" gorm:"type:timestamp;"`         //登录时间
	CreateBy      string    `json:"createBy" gorm:"type:varchar(128);"`       //创建人
	UpdateBy      string    `json:"updateBy" gorm:"type:varchar(128);"`       //更新者
	DataScope     string    `json:"dataScope" gorm:"-"`                       //数据
	Params        string    `json:"params" gorm:"-"`                          //
	Remark        string    `json:"remark" gorm:"type:varchar(255);"`         //备注
	Msg           string    `json:"msg" gorm:"type:varchar(255);"`
	BaseModel
}

func (LoginLog) TableName() string {
	return "prince_loginlog"
}

func (e *LoginLog) Create() (LoginLog, error) {
	var doc LoginLog
	e.CreateBy = "0"
	e.UpdateBy = "0"
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}
