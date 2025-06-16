package models

import (
	"gorm.io/gorm"
	"time"
)

type AddressBook struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey; comment:地址唯一标识"`
	UserId        uint64 `gorm:"comment '用户id'"`
	Consignee     string `gorm:"type:varchar(60) comment '收货人'"`
	Gender        uint32 `gorm:"comment '收货人性别 1男 0女'"`
	Phone         string `gorm:"type:varchar(11) comment '手机号'"`
	StressAddress string `gorm:"type:varchar(255) comment '详细地址'"`
	City          string `gorm:"type:varchar(60) comment '城市'"`
	State         string `gorm:"type:varchar(60) comment '省份/州'"`
	Country       string `gorm:"type:varchar(60) comment '国家'"`
	ZipCode       string `gorm:"type:char(6) comment '邮编'"`
	Label         string `gorm:"comment '标签'"`
	IsDefault     bool   `gorm:"comment '是否为默认地址'"`
}

func (a AddressBook) GetID() uint64 {
	return a.ID
}

func (a AddressBook) GetCreatedTime() time.Time {
	return a.CreatedAt
}
