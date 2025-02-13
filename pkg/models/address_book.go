package models

import "gorm.io/gorm"

type AddressBook struct {
	gorm.Model
	UserId        uint32 `gorm:"comment '用户id'"`
	Consignee     string `gorm:"type:varchar(60) comment '收货人'"`
	Gender        uint32 `gorm:"comment '收货人性别 1男 0女'"`
	Phone         string `gorm:"type:varchar(11) unique comment '手机号'"`
	StressAddress string `gorm:"type:varchar(255) comment '详细地址'"`
	City          string `gorm:"type:varchar(60) comment '城市'"`
	State         string `gorm:"type:varchar(60) comment '省份/州'"`
	Country       string `gorm:"type:varchar(60) comment '国家'"`
	ZipCode       int32  `gorm:"type:varchar(6) comment '邮编'"`
	Label         string `gorm:"comment '标签'"`
	IsDefault     bool   `gorm:"comment '是否为默认地址'"`
}
