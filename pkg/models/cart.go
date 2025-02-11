package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint32 `gorm:"comment '用户ID'"`
	ProductID uint32 `gorm:"comment '商品ID'"`
	Num       int    `gorm:"comment '商品数量'"`
}
