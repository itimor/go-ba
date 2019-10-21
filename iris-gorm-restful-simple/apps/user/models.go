package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);not null;unique"`
	Password string `gorm:"type:varchar(100);not null"`
	Roles    []Role `gorm:"many2many:user_roles;"`
	Avatar   string `gorm:"not null"`
	IsActive bool   `gorm:"default:true"`
	IsAdmin  bool   `gorm:"default:true"`
}

type Role struct {
	gorm.Model
	name string `gorm:"type:varchar(100);not null;unique"`
	Desc string `sql:"type:text"`
}
