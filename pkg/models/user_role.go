package models

import (
	"gorm.io/gorm"
	"time"
)

type UserRole struct {
	gorm.Model
	ID     uint64 `gorm:"primaryKey;"` //唯一标识
	UserID uint64 `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleID uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`
}

func (u UserRole) GetID() uint64 {
	return u.ID
}

func (u UserRole) GetCreatedTime() time.Time {
	return u.CreatedAt
}
