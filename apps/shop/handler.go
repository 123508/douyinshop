package main

import (
	"context"
	"github.com/123508/douyinshop/pkg/els"
	"github.com/123508/douyinshop/pkg/errorno"
	"log"
	"strconv"
	"strings"

	ba "github.com/123508/douyinshop/kitex_gen/product"
	pb "github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/models"
	"go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

var RepeatedShop = &errorno.BasicMessageError{Code: 404, Message: "你已注册店铺，请勿重复注册"}

var ShopNotFound = &errorno.BasicMessageError{Code: 404, Message: "无法找到店铺"}

var ProductLoss = &errorno.BasicMessageError{Code: 404, Message: "商品信息丢失"}

var FailFetchProductList = &errorno.BasicMessageError{Code: 404, Message: "无法获取商品列表"}

type ShopServiceImpl struct {
	db         *gorm.DB
	etcdClient *clientv3.Client
	leaseID    clientv3.LeaseID
}

// Register 注册店铺
func (s *ShopServiceImpl) Register(ctx context.Context, req *pb.RegisterShopReq) (*pb.RegisterShopResp, error) {
	if err := s.db.Model(&models.Shop{}).Where("user_id = ?", req.UserId).First(&models.Shop{}).Error; err == nil {
		return nil, RepeatedShop
	}
	shop := models.Shop{
		UserId:      req.UserId,
		Name:        req.ShopName,
		Address:     req.ShopAddress,
		Description: req.ShopDescription,
		Avatar:      req.ShopAvatar,
	}
	result := s.db.Create(&shop)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pb.RegisterShopResp{ShopId: shop.Id}, nil
}

// GetShopId 获取用户所开的店铺id
func (s *ShopServiceImpl) GetShopId(ctx context.Context, req *pb.GetShopIdReq) (*pb.GetShopIdResp, error) {
	var shop models.Shop
	result := s.db.Where("user_id = ?", req.UserId).First(&shop)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &pb.GetShopIdResp{ShopId: 0}, nil
		}
		return nil, result.Error
	}
	return &pb.GetShopIdResp{ShopId: shop.Id}, nil
}

// GetShopInfo 获取店铺信息
func (s *ShopServiceImpl) GetShopInfo(ctx context.Context, req *pb.GetShopInfoReq) (*pb.GetShopInfoResp, error) {
	var shop models.Shop
	result := s.db.Where("id = ?", req.ShopId).First(&shop)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ShopNotFound
		}
		return nil, result.Error
	}
	return &pb.GetShopInfoResp{
		ShopName:        shop.Name,
		ShopAddress:     shop.Address,
		ShopDescription: shop.Description,
		ShopAvatar:      shop.Avatar,
	}, nil
}

// UpdateShopInfo 更新店铺信息
func (s *ShopServiceImpl) UpdateShopInfo(ctx context.Context, req *pb.UpdateShopInfoReq) (*pb.UpdateShopInfoResp, error) {
	var shop models.Shop
	result := s.db.Where("id = ?", req.ShopId).First(&shop)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &pb.UpdateShopInfoResp{Res: false}, nil
		}
		return nil, result.Error
	}
	shop.Name = req.ShopName
	shop.Address = req.ShopAddress
	shop.Description = req.ShopDescription
	shop.Avatar = req.ShopAvatar
	updateResult := s.db.Save(&shop)
	if updateResult.Error != nil {
		return &pb.UpdateShopInfoResp{Res: false}, updateResult.Error
	}
	return &pb.UpdateShopInfoResp{Res: true}, nil
}

// AddProduct 添加商品
func (s *ShopServiceImpl) AddProduct(ctx context.Context, req *pb.AddProductReq) (*pb.AddProductResp, error) {
	if req.Product == nil {
		return nil, ProductLoss
	}
	categoriesStr := strings.Join(req.Product.Categories, ",")
	product := models.Product{
		ShopId:      uint(req.ShopId),
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Picture:     req.Product.Picture,
		Price:       float32(req.Product.Price),
		Categories:  categoriesStr,
		Status:      req.Status, // 修正状态字段来源
	}
	result := s.db.Create(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	// 更新ES中的商品信息
	err := els.UpdateProduct(
		&ba.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Categories:  req.Product.Categories,
			Sales:       0,
		},
	)
	if err != nil {
		return nil, err
	}

	return &pb.AddProductResp{ProductId: product.Id}, nil
}

// DeleteProduct 删除商品
func (s *ShopServiceImpl) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.DeleteProductResp, error) {
	var product models.Product
	result := s.db.Where("id = ? AND shop_id = ?", req.ProductId, req.ShopId).Delete(&product)
	if result.Error != nil {
		return &pb.DeleteProductResp{Res: false}, result.Error
	}

	// 删除ES中的商品信息
	err := els.DeleteProduct(req.ProductId)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResp{Res: result.RowsAffected > 0}, nil
}

// UpdateProduct 更新商品
func (s *ShopServiceImpl) UpdateProduct(ctx context.Context, req *pb.UpdateProductReq) (*pb.UpdateProductResp, error) {
	var product models.Product
	result := s.db.Where("id = ? AND shop_id = ?", req.Product.Id, req.ShopId).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return &pb.UpdateProductResp{Res: false}, nil
		}
		return nil, result.Error
	}
	// 将 req.Product.Categories 从 []string 转换为 string
	categoriesStr := strings.Join(req.Product.Categories, ",")
	product.Name = req.Product.Name
	product.Description = req.Product.Description
	product.Picture = req.Product.Picture
	product.Price = float32(req.Product.Price)
	product.Categories = categoriesStr
	product.Status = req.Product.Status
	updateResult := s.db.Save(&product)
	if updateResult.Error != nil {
		return &pb.UpdateProductResp{Res: false}, updateResult.Error
	}

	// 更新ES中的商品信息
	err := els.UpdateProduct(req.Product)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProductResp{Res: true}, nil
}

// GetProductList 获取商品列表
func (s *ShopServiceImpl) GetProductList(ctx context.Context, req *pb.GetProductListReq) (*pb.GetProductListResp, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	offset := (req.Page - 1) * req.PageSize

	log.Printf("Received GetProductList request for shop ID: %d, page: %d, page size: %d", req.ShopId, req.Page, req.PageSize)

	var products []models.Product
	result := s.db.Where("shop_id = ?", req.ShopId).
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find(&products)

	if result.Error != nil {
		log.Printf("Error fetching product list for shop ID %d: %v", req.ShopId, result.Error)
		return nil, FailFetchProductList
	}
	if result.RowsAffected == 0 {
		log.Printf("No products found for shop ID %d on page %d with page size %d", req.ShopId, req.Page, req.PageSize)
		return &pb.GetProductListResp{Products: []*ba.Product{}}, nil
	}

	pbProducts := make([]*ba.Product, len(products))
	for i, p := range products {
		category := make([]string, 0)
		for _, categoryIdStr := range strings.Split(p.Categories, ",") {
			var categoryResult models.Category
			categoryId, _ := strconv.Atoi(strings.Trim(categoryIdStr, " "))
			s.db.First(&categoryResult, categoryId)
			category = append(category, categoryResult.Name)
		}
		pbProducts[i] = &ba.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       float32(p.Price),
			Categories:  category,
		}
	}

	log.Printf("Successfully fetched %d products for shop ID %d", len(pbProducts), req.ShopId)
	return &pb.GetProductListResp{Products: pbProducts}, nil
}
