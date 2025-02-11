package models

import (
	"gorm.io/gorm"
)

type Shop struct {
	gorm.Model
	Id          uint32 `gorm:"primary_key;comment '店铺ID'"`
	Name        string `gorm:"type:varchar(60) comment '店铺名称'"`
	Description string `gorm:"comment '店铺描述'"`
	Avatar      string `gorm:"type varchar(200) comment '店铺头像'"`
	Address     string `gorm:"type varchar(200) comment '店铺地址'"`
	UserId      uint32 `gorm:"comment '店主ID'"`
}
