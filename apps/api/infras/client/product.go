package client

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/kitex_gen/product"
	"github.com/123508/douyinshop/kitex_gen/product/productcatalogservice"
	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var productClient productcatalogservice.Client

func initProductRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := productcatalogservice.NewClient(
		config.Conf.ProductConfig.ServiceName,             // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	productClient = c
}

type ProductItem struct {
	Id           uint32   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Picture      string   `json:"picture"`
	Price        string   `json:"price"`
	CategoryName []string `json:"categories"`
	Sales        int      `json:"sales"`
	ShopId       uint32   `json:"shop_id"`
}

// ListProducts 获取商品列表
// page 页码
// size 每页数量
// category 分类
// 返回商品列表
func ListProducts(ctx context.Context, page int, size int, category string) ([]ProductItem, error) {
	req := &product.ListProductsReq{
		Page:         int32(page),
		PageSize:     int64(size),
		CategoryName: category,
	}
	resp, err := productClient.ListProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	var products []ProductItem
	for _, p := range resp.Products {
		products = append(products, ProductItem{
			Id:           p.Id,
			Name:         p.Name,
			Price:        fmt.Sprintf("%.2f", p.Price),
			Description:  p.Description,
			CategoryName: p.Categories,
			Picture:      p.Picture,
			ShopId:       p.ShopId,
		})
	}
	return products, nil
}

// GetProductDetail 获取商品详情
// productId 商品ID
// 返回商品详情
func GetProductDetail(ctx context.Context, productId int) (*ProductItem, error) {
	req := &product.GetProductReq{
		Id: uint32(productId),
	}
	resp, err := productClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	result := &ProductItem{
		Id:           resp.Product.Id,
		Name:         resp.Product.Name,
		Price:        fmt.Sprintf("%.2f", resp.Product.Price),
		Description:  resp.Product.Description,
		CategoryName: resp.Product.Categories,
		Picture:      resp.Product.Picture,
		Sales:        int(resp.Product.Sales),
		ShopId:       resp.Product.ShopId,
	}
	return result, nil
}

// SearchProducts 搜索商品
// keyword 关键字
// page 页码
// pageSize 每页数量
// 返回商品列表
func SearchProducts(ctx context.Context, keyword string, page int, pageSize int) ([]ProductItem, error) {
	req := &product.SearchProductsReq{
		Query:    keyword,
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}
	resp, err := productClient.SearchProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	var products []ProductItem
	for _, p := range resp.Results {
		products = append(products, ProductItem{
			Id:           p.Id,
			Name:         p.Name,
			Price:        fmt.Sprintf("%.2f", p.Price),
			Description:  p.Description,
			CategoryName: p.Categories,
			Picture:      p.Picture,
			ShopId:       p.ShopId,
		})
	}
	return products, nil
}
