package els

import (
	"github.com/123508/douyinshop/kitex_gen/product"
	"testing"
)

func TestSearchProduct(t *testing.T) {
	name := "手机"
	got, err := SearchProduct(name, 1, 2)
	if err != nil {
		t.Errorf("SearchProduct() error = %v", err)
		return
	}
	t.Logf("SearchProduct() got = %v", got)
}

func TestUpdateProduct(t *testing.T) {
	var productItem = product.Product{
		Id:          7,
		Name:        "商品名称2",
		Description: "商品描述2",
		Picture:     "商品图片2",
		Price:       100,
		Categories:  []string{"商品类别1", "商品类别2"},
		Sales:       100,
	}
	err := UpdateProduct(&productItem)
	if err != nil {
		t.Errorf("UpdateProduct() error = %v", err)
		return
	}
	t.Logf("UpdateProduct() success")
}
