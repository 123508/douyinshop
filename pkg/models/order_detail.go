package models

import "gorm.io/gorm"

type OrderDetail struct {
	gorm.Model
	Name      string  `gorm:"comment '商品名称'"`
	Image     string  `gorm:"comment '商品图片'"`
	OrderId   int     `gorm:"comment '订单id'"`
	ProductId int     `gorm:"comment '商品id'"`
	Number    int     `gorm:"comment '商品数量'"`
	Amount    float64 `gorm:"type:decimal(10,2) comment '支付金额,精确到后两位'"`
	Status    int     `gorm:"comment '订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单'"`
}
