package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderStatusLog struct {
	gorm.Model
	ID          uint64     `gorm:"primary_key;comment '订单日志ID'"`
	OrderId     uint64     `gorm:"comment '订单详情id'"`
	Status      uint32     `gorm:"comment '订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单'"`
	StartTime   *time.Time `gorm:"comment '状态开始时间'"`
	EndTime     *time.Time `gorm:"comment '状态结束时间'"`
	Description string     `gorm:"comment '状态描述'"`
	Version     uint32     `gorm:"comment '当前版本号'"`
}

func (o OrderStatusLog) GetID() uint64 {
	return o.ID
}

func (o OrderStatusLog) GetCreatedTime() time.Time {
	return o.CreatedAt
}
