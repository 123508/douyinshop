package models

import (
	"gorm.io/gorm"
	"time"
)

type Permission struct {
	gorm.Model
	ID             uint64 `gorm:"primaryKey;column:id" json:"id"`                                  // 唯一标识
	Code           string `gorm:"uniqueIndex;size:100;not null;column:code" json:"code"`           // 权限唯一编码（如ORDER_VIEW等，便于代码调用）
	PermissionName string `gorm:"size:100;not null;column:permission_name" json:"permission_name"` // 权限名称
	Description    string `gorm:"size:100;column:description" json:"description"`                  // 权限详细描述
	ParentID       uint64 `gorm:"column:parent_id" json:"parent_id"`                               // 父权限ID，支持多级权限结构（如菜单、按钮）
	Type           string `gorm:"size:20;column:type" json:"type"`                                 // 权限类型（菜单、按钮、接口等）
	Resource       string `gorm:"size:200;column:resource" json:"resource"`                        // 资源标识（如API路径、菜单路由）
	Method         string `gorm:"size:10;column:method" json:"method"`                             // 操作方法（如GET/POST/PUT/DELETE等，API权限控制时用）
	Status         int    `gorm:"default:1;column:status" json:"status"`                           // 状态（启用/禁用）
}

func (p Permission) GetID() uint64 {
	return p.ID
}

func (p Permission) GetCreatedTime() time.Time {
	return p.CreatedAt
}
