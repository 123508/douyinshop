package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID         uint64 `gorm:"primary_key;comment '分类ID'"`
	Name       string `gorm:"type:varchar(32) comment '分类名称'"`
	Status     bool   `gorm:"comment '分类是否启用'"`
	CreateUser uint64 `gorm:"comment '创建人id'"`
	UpdateUser uint64 `gorm:"comment '更新人id'"`
}

func (c Category) GetID() uint64 {
	return c.ID
}
