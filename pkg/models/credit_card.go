package models

import (
	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model
	CreditCardNumber          string `gorm:"type:varchar(16) comment '信用卡号'"`
	CreditCardCvv             int32  `gorm:"type:varchar(3) comment '信用卡CVV'"`
	CreditCardExpirationYear  int32  `gorm:"comment '信用卡过期年份'"`
	CreditCardExpirationMonth int32  `gorm:"comment '信用卡过期月份'"`
}
