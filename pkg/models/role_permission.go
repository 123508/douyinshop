package models

import (
	"gorm.io/gorm"
	"time"
)

type RolePermission struct {
	gorm.Model
	RoleID       uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`             //角色ID
	PermissionID uint64 `gorm:"primaryKey;column:permission_id" json:"permission_id"` //权限ID
	CreatedBy    uint64 `gorm:"column:created_by" json:"created_by"`                  //分配人ID
}

func (r RolePermission) GetID() uint64 {
	return r.RoleID
}

func (r RolePermission) GetCreatedTime() time.Time {
	return r.CreatedAt
}
