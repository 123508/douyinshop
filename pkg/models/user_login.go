package models

import (
	"gorm.io/gorm"
)

type UserLogin struct {
	gorm.Model
	UserId   uint32 `gorm:"unique comment '用户ID'"`
	Password string `gorm:"size:200;notnull comment '加密后的密码'"`
}
