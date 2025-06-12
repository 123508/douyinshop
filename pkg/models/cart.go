package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID        uint64 `gorm:"primary_key;comment '购物明细ID'"`
	UserId    uint64 `gorm:"comment '用户ID'"`
	ProductID uint64 `gorm:"comment '商品ID'"`
	Num       int    `gorm:"comment '商品数量'"`
	Status    int    `gorm:"comment '状态'"` // 0: 选中, 1: 未选中, 2: 商品下架

}

func (c Cart) GetID() uint64 {
	return c.ID
}
