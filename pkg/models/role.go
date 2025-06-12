package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint64 `gorm:"primaryKey;column:id" json:"id"`                       // 唯一标识
	Code        string `gorm:"uniqueIndex;size:50;not null;column:code" json:"code"` // 角色唯一编码
	RoleName    string `gorm:"size:50;not null;column:role_name" json:"role_name"`   // 角色名称
	Description string `gorm:"size:100;column:description" json:"description"`       // 角色描述
	Status      int    `gorm:"default:1;column:status" json:"status"`                // 1:启用 0:禁用
	CreatedBy   uint64 `gorm:"column:created_by" json:"created_by"`                  // 创建人用户ID
}

func (r Role) GetID() uint64 {
	return r.ID
}
