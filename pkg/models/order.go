package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserId          uint32    `gorm:"comment '用户Id'"`
	Number          string    `gorm:"type:varchar(50) comment '订单号'"`
	Status          int       `gorm:"comment '订单状态 1待付款 2待接单 3已接单 4派送中 5待收货 6已完成 7已取消 8退款'"`
	AddressBookId   int       `gorm:"comment '地址id'"`
	CheckOutTime    time.Time `gorm:"comment '结账时间'"`
	PayMethod       int       `gorm:"comment '支付方式 1微信 2支付宝 3银行卡'"`
	PayStatus       int       `gorm:"comment '支付状态 0未支付 1已支付 2已取消'"`
	Amount          float64   `gorm:"type:decimal(10,2) comment '支付金额,精确到后两位'"`
	Remark          string    `gorm:"comment '用户备注'"`
	Phone           string    `gorm:"type:varchar(11) comment '用户电话'"`
	Address         string    `gorm:"comment '下单地址'"`
	UserName        string    `gorm:"type:varchar(60) comment '用户名称'"`
	Consignee       string    `gorm:"type:varchar(60) comment '收货人'"`
	CancelReason    string    `gorm:"comment '取消订单原因'"`
	RejectionReason string    `gorm:"comment '商家拒单原因'"`
	CancelTime      time.Time `gorm:"comment '取消订单时间'"`
}
