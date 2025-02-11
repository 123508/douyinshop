package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id          uint32     `gorm:"primary_key;comment '商品ID'"`
	ShopId      uint       `gorm:"comment '所属于店铺ID'"`
	Name        string     `gorm:"comment '商品名称'"`
	Description string     `gorm:"comment '商品描述'"`
	Picture     string     `gorm:"comment '商品图片'"`
	Price       float32    `gorm:"type:decimal(10,2) comment '商品价格,精确到后两位'"`
	Categories  []Category `gorm:"many2many:product_categories;comment '分类名称'"`
	Status      bool       `gorm:"comment '商品是否出售'"`
}
