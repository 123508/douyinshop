package models

import (
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	gorm.Model
	ID                        uint64 `gorm:"primary_key;comment '信用卡ID'"`
	CreditCardNumber          string `gorm:"type:varchar(16) comment '信用卡号'"`
	CreditCardCvv             int32  `gorm:"type:varchar(3) comment '信用卡CVV'"`
	CreditCardExpirationYear  int32  `gorm:"comment '信用卡过期年份'"`
	CreditCardExpirationMonth int32  `gorm:"comment '信用卡过期月份'"`
}

func (c CreditCard) GetID() uint64 {
	return c.ID
}

func (c CreditCard) GetCreatedTime() time.Time {
	return c.CreatedAt
}
