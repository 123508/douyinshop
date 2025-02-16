package main

import (
	"context"
	product "github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/pkg/models"
	"strconv"
	"strings"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts implements the ProductCatalogServiceImpl interface.
// 若分类名为空，则返回所以商品中的第page页的pageSize个商品
// 当商品不存在时，返回空列表
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	var products []models.Product
	page := int(req.Page)
	pageSize := int(req.PageSize)
	productList := make([]*product.Product, 0)
	if req.CategoryName == "" {
		database.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
		for _, productItem := range products {
			category := make([]string, 0)
			for _, categoryIdStr := range strings.Split(productItem.Categories, ",") {
				var categoryResult models.Category
				categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
				database.First(&categoryResult, categoryId)
				category = append(category, categoryResult.Name)
			}
			productList = append(productList, &product.Product{
				Id:          uint32(productItem.ID),
				Name:        productItem.Name,
				Description: productItem.Description,
				Picture:     productItem.Picture,
				Price:       productItem.Price,
				Categories:  category,
			})
		}
	}
	resp = &product.ListProductsResp{
		Products: productList,
	}
	return
}

// GetProduct implements the ProductCatalogServiceImpl interface.
// 当商品不存在时，返回nil
// 当商品存在时，返回商品信息
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	result := models.Product{}
	database.First(&result, req.Id)
	if result.Id == 0 {
		resp = &product.GetProductResp{
			Product: nil,
		}
	} else {
		category := make([]string, 0)
		for _, categoryIdStr := range strings.Split(result.Categories, ",") {
			var categoryResult models.Category
			categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
			database.First(&categoryResult, categoryId)
			category = append(category, categoryResult.Name)
		}
		resp = &product.GetProductResp{
			Product: &product.Product{
				Id:          uint32(result.ID),
				Name:        result.Name,
				Description: result.Description,
				Picture:     result.Picture,
				Price:       result.Price,
				Categories:  category,
			},
		}
	}
	return
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// TODO: Your code here...
	return
}
