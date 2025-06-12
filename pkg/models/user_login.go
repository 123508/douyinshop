package models

import (
	"gorm.io/gorm"
)

type UserLogin struct {
	gorm.Model
	ID       uint64 `gorm:"primary_key;comment '用户登录表ID'"`
	UserId   uint64 `gorm:"unique comment '用户ID'"`
	Password string `gorm:"size:200;notnull comment '加密后的密码'"`
}

func (u UserLogin) GetID() uint64 {
	return u.ID
}
