package models

import "gorm.io/gorm"

type OrderDetail struct {
	gorm.Model
	Name      string  `gorm:"comment '商品名称'"`
	Image     string  `gorm:"comment '商品图片'"`
	OrderId   uint32  `gorm:"comment '订单id'"`
	ProductId uint32  `gorm:"comment '商品id'"`
	Number    uint32  `gorm:"comment '商品数量'"`
	Amount    float32 `gorm:"type:decimal(10,2) comment '支付金额,精确到后两位'"`
}
