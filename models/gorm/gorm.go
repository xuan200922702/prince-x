package gorm

import (
	"github.com/jinzhu/gorm"
	"prince-x/models"
)

func AutoMigrate(db *gorm.DB) error {
	db.SingularTable(true)
	return db.AutoMigrate(
		new(models.PrinceUser),
		new(models.PrinceRole),
		new(models.LoginLog),
		new(models.CasbinRule),
		new(models.Menu),
		new(models.RoleMenu),
		new(models.Dept),
		new(models.PrinceRoleDept),
	).Error
}
