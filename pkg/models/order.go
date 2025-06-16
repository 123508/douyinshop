package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	gorm.Model
	ID            uint64  `gorm:"primarykey"`
	UserId        uint64  `gorm:"comment '用户Id'"`
	ShopId        uint64  `gorm:"comment  '商家Id'"`
	Number        string  `gorm:"type:varchar(50) comment '订单号'"`
	AddressBookId uint64  `gorm:"comment '地址id'"`
	PayMethod     int32   `gorm:"comment '支付方式 1微信 2支付宝 3银行卡'"`
	Amount        float32 `gorm:"type:decimal(10,2) comment '支付金额,精确到后两位'"`
	Remark        string  `gorm:"comment '用户备注'"`
	Phone         string  `gorm:"type:varchar(11) comment '用户电话'"`
	Address       string  `gorm:"comment '下单地址'"`
	UserName      string  `gorm:"type:varchar(60) comment '用户名称'"`
	Consignee     string  `gorm:"type:varchar(60) comment '收货人'"`
	TransactionId string  `gorm:"type:varchar(50) comment '支付交易号'"`
	FinalStatus   uint32  `gorm:"comment '订单最终状态'"`
	FinalVersion  uint32  `gorm:"comment '最终版本号'"`
}

func (o Order) GetID() uint64 {
	return o.ID
}

func (o Order) GetCreatedTime() time.Time {
	return o.CreatedAt
}
