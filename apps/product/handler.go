package main

import (
	"context"
	"encoding/json"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/pkg/els"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Category = "category"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

func (s *ProductCatalogServiceImpl) GetCategoryFromProduct(ctx context.Context, Categories string, productId uint64) ([]models.Category, error) {

	simple := util.SimpleCacheComponent[uint64, []models.Category]{
		Rds:       Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(Category, productId),
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() ([]models.Category, error) {

			categoryIdList := make([]uint64, 0)
			for _, categoryIdStr := range strings.Split(Categories, ",") {
				categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
				categoryIdList = append(categoryIdList, uint64(categoryId))
			}

			category := make([]models.Category, 0)

			if err := database.Model(&models.Category{}).Where("id IN ?", categoryIdList).Find(&category).Error; err != nil {
				return nil, err
			}
			return category, nil
		},
		Expires: time.Duration(rand.Intn(3)+3) * time.Second,
	}

	return simple.QueryWithCache()
}

// ListProducts implements the ProductCatalogServiceImpl interface.
// 获取商品列表接口
// 若分类名为空，则返回所以商品中的第page页的pageSize个商品
// 若分类名不为空，则返回指定分类名的第page页的pageSize个商品
// 当商品不存在时，返回空列表
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {

	//对请求页数和页长进行判断
	if req.Page < 1 || req.PageSize < 1 {
		return nil, BadPageOrPageSize
	}

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
				Id:          productItem.ID,
				Name:        productItem.Name,
				Description: productItem.Description,
				Picture:     productItem.Picture,
				Price:       productItem.Price,
				Categories:  category,
				Sales:       productItem.Sales,
				ShopId:      productItem.ShopId,
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
				Id:          productItem.ID,
				Name:        productItem.Name,
				Description: productItem.Description,
				Picture:     productItem.Picture,
				Price:       productItem.Price,
				Categories:  category,
				Sales:       productItem.Sales,
				ShopId:      productItem.ShopId,
			})
		}
	}
	resp = &product.ListProductsResp{
		Products: productList,
	}
	return
}

// GetProduct implements the ProductCatalogServiceImpl interface.
// 查找指定商品接口
// 当商品不存在时，返回nil
// 当商品存在时，返回商品信息
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	result := models.Product{}
	database.First(&result, req.Id)
	if result.ID == 0 {
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
				Id:          result.ID,
				Name:        result.Name,
				Description: result.Description,
				Picture:     result.Picture,
				Price:       result.Price,
				Categories:  category,
				Sales:       result.Sales,
				ShopId:      result.ShopId,
			},
		}
	}
	return
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
// 搜索商品接口
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
			Id:          productResult.ID,
			Name:        productResult.Name,
			Description: productResult.Description,
			Picture:     productResult.Picture,
			Price:       productResult.Price,
			Categories:  category,
			Sales:       productResult.Sales,
			ShopId:      productResult.ShopId,
		})
	}
	resp = &product.SearchProductsResp{
		Results: productList,
	}
	return
}
