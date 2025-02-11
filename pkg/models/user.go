package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id     uint32 `gorm:"primary_key;comment '用户ID'"`
	Name   string `gorm:"type:varchar(60) comment '用户名称'"`
	Email  string `gorm:"unique;comment '用户邮箱'"`
	Phone  string `gorm:"type:varchar(11) unique comment '用户手机号'"`
	Gender byte   `gorm:"comment '用户性别 1男 0女'"`
	Avatar string `gorm:"type varchar(200) comment '用户头像'"`
}
