package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(64);uniqueIndex;not null;comment:用户名"`
	Password string `gorm:"type:varchar(128);not null;comment:密码"`
	Name     string `gorm:"type:varchar(64);comment:姓名"`
	Email    string `gorm:"type:varchar(128);index;comment:邮箱"`
	Status   int    `gorm:"type:tinyint;default:1;comment:状态 1正常 0禁用"`
}

func (User) TableName() string {
	return "users"
}
