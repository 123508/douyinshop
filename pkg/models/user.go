package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     uint32 `gorm:"primary_key;auto_increment;comment '用户ID'"`
	Name   string `gorm:"type:varchar(60) comment '用户名称'"`
	Email  string `gorm:"unique;comment '用户邮箱'"`
	Phone  string `gorm:"type:varchar(11) unique comment '用户手机号'"`
	Gender uint32 `gorm:"comment '用户性别 1男 0女'"`
	Avatar string `gorm:"type varchar(200) comment '用户头像'"`
}
