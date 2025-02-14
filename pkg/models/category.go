package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name       string `gorm:"type:varchar(32) comment '分类名称'"`
	Status     bool   `gorm:"comment '分类是否启用'"`
	CreateUser uint32 `gorm:"comment '创建人id'"`
	UpdateUser uint32 `gorm:"comment '更新人id'"`
}
