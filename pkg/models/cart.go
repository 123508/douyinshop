package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint32 `gorm:"comment '用户ID'"`
	ProductID uint32 `gorm:"comment '商品ID'"`
	Num       int    `gorm:"comment '商品数量'"`
	Status    int    `gorm:"comment '状态'"` // 0: 选中, 1: 未选中, 2: 商品下架
}
