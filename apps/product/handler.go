package main

import (
	"context"
	product "github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/pkg/els"
	"github.com/123508/douyinshop/pkg/models"
	"strconv"
	"strings"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// ListProducts implements the ProductCatalogServiceImpl interface.
// 若分类名为空，则返回所以商品中的第page页的pageSize个商品
// 若分类名不为空，则返回指定分类名的第page页的pageSize个商品
// 当商品不存在时，返回空列表
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	var products []models.Product
	page := int(req.Page)
	pageSize := int(req.PageSize)
	productList := make([]*product.Product, 0)
	if req.CategoryName == "" { // 分类名为空时，返回所有类型商品
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
				Id:          productItem.Id,
				Name:        productItem.Name,
				Description: productItem.Description,
				Picture:     productItem.Picture,
				Price:       productItem.Price,
				Categories:  category,
				Sales:       productItem.Sales,
				ShopId:      uint32(productItem.ShopId),
			})
		}
	} else { // 分类名不为空时，返回指定类型商品
		database.Offset((page-1)*pageSize).Limit(pageSize).Where("categories like ?", "%"+req.CategoryName+"%").Find(&products)
		for _, productItem := range products {
			category := make([]string, 0)
			for _, categoryIdStr := range strings.Split(productItem.Categories, ",") {
				var categoryResult models.Category
				categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
				database.First(&categoryResult, categoryId)
				category = append(category, categoryResult.Name)
			}
			productList = append(productList, &product.Product{
				Id:          productItem.Id,
				Name:        productItem.Name,
				Description: productItem.Description,
				Picture:     productItem.Picture,
				Price:       productItem.Price,
				Categories:  category,
				Sales:       productItem.Sales,
				ShopId:      uint32(productItem.ShopId),
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
				Id:          result.Id,
				Name:        result.Name,
				Description: result.Description,
				Picture:     result.Picture,
				Price:       result.Price,
				Categories:  category,
				Sales:       result.Sales,
				ShopId:      uint32(result.ShopId),
			},
		}
	}
	return
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
// 当搜索结果为空时，返回空列表
// 当搜索结果不为空时，返回搜索结果
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	result, err := els.SearchProduct(req.Query, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	productList := make([]*product.Product, 0)
	for _, productId := range result {
		var productResult models.Product
		database.First(&productResult, productId)
		category := make([]string, 0)
		for _, categoryIdStr := range strings.Split(productResult.Categories, ",") {
			var categoryResult models.Category
			categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
			database.First(&categoryResult, categoryId)
			category = append(category, categoryResult.Name)
		}
		productList = append(productList, &product.Product{
			Id:          productResult.Id,
			Name:        productResult.Name,
			Description: productResult.Description,
			Picture:     productResult.Picture,
			Price:       productResult.Price,
			Categories:  category,
			Sales:       productResult.Sales,
			ShopId:      uint32(productResult.ShopId),
		})
	}
	resp = &product.SearchProductsResp{
		Results: productList,
	}
	return
}
