package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderDetail struct {
	gorm.Model
	ID        uint64  `gorm:"primary_key;comment '订单明细ID'"`
	Name      string  `gorm:"comment '商品名称'"`
	Image     string  `gorm:"comment '商品图片'"`
	OrderId   uint64  `gorm:"comment '订单id'"`
	ProductId uint64  `gorm:"comment '商品id'"`
	Number    uint32  `gorm:"comment '商品数量'"`
	Amount    float32 `gorm:"type:decimal(10,2) comment '支付金额,精确到后两位'"`
}

func (o OrderDetail) GetID() uint64 {
	return o.ID
}

func (o OrderDetail) GetCreatedTime() time.Time {
	return o.CreatedAt
}
