package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/123508/douyinshop/pkg/els"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"

	ba "github.com/123508/douyinshop/kitex_gen/product"
	pb "github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/models"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	serviceName = "shop"
	product     = "product"
)

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

var RepeatedShop = &errorno.BasicMessageError{Code: 404, Message: "你已注册店铺，请勿重复注册"}

var ShopNotFound = &errorno.BasicMessageError{Code: 404, Message: "无法找到店铺"}

var ProductLoss = &errorno.BasicMessageError{Code: 404, Message: "商品信息丢失"}

var FailFetchProductList = &errorno.BasicMessageError{Code: 404, Message: "无法获取商品列表"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}

type ShopServiceImpl struct {
	db         *gorm.DB
	etcdClient *clientv3.Client
	leaseID    clientv3.LeaseID
}

func DelShopProductCache(ctx context.Context, rds *redis.Client, shopId uint64) {
	// 构建 pattern
	pattern := util.TakeKey(product, shopId, "*")
	var cursor uint64 = 0
	var batchSize int64 = 100 // 可调整

	for {
		keys, nextCursor, err := rds.Scan(ctx, cursor, pattern, batchSize).Result()
		if err != nil {
			util.LogError("redis扫描错误", "DelShopProductCache", "", err)
			return
		}
		if len(keys) > 0 {
			pipe := rds.Pipeline()
			for _, key := range keys {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				util.LogError("redis pipeline删除错误", "DelShopProductCache", "", err)
				return
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
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
	return &pb.RegisterShopResp{ShopId: shop.ID}, nil
}

// GetShopId 获取用户所开的店铺id
func (s *ShopServiceImpl) GetShopId(ctx context.Context, req *pb.GetShopIdReq) (*pb.GetShopIdResp, error) {

	var shop models.Shop

	key := util.TakeKey(serviceName, "shopId", req.UserId)
	shopId, err := rds.Get(ctx, key).Result()

	if err != nil {
		util.LogError("查询shopId缓存失败", "GetShopId", "", err)
	}

	id, err := strconv.Atoi(shopId)

	//查询缓存失败
	if err != nil {

		result := s.db.Where("user_id = ?", ctx.Value("userId").(uint32)).First(&shop)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return &pb.GetShopIdResp{ShopId: 0}, nil
			}
			return nil, result.Error
		}

		//存入缓存
		err := rds.Set(ctx, key, util.TakeKey(shop.ID), time.Duration(rand.Intn(15)+30)*time.Minute).Err()

		if err != nil {
			util.LogError("缓存数据失败", "GetShopId", "", err)
		}

	} else {
		shop.ID = uint64(id)
		log.WithFields(log.Fields{
			"方法名": "GetShopId",
		}).Info("查询缓存成功")

	}

	return &pb.GetShopIdResp{ShopId: shop.ID}, nil
}

// GetShopInfo 获取店铺信息
func (s *ShopServiceImpl) GetShopInfo(ctx context.Context, req *pb.GetShopInfoReq) (*pb.GetShopInfoResp, error) {

	var shop models.Shop

	key := util.TakeKey(serviceName, "shopInfo", req.ShopId)

	jsonData, err := rds.Get(ctx, key).Result()

	if err != nil {
		util.LogError("查询shopInfo缓存失败", "GetShopInfo", "", err)
	}

	err = json.Unmarshal([]byte(jsonData), &shop)

	//查询缓存失败
	if err != nil {

		result := s.db.Where("id = ?", req.ShopId).First(&shop)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return nil, ShopNotFound
			}
			return nil, result.Error
		}

		//放入缓存
		jsonData, _ := json.Marshal(&shop)

		err := rds.Set(ctx, key, string(jsonData), time.Duration(rand.Intn(15)+30)*time.Minute).Err()

		if err != nil {
			util.LogError("缓存数据失败", "GetShopInfo", "", err)
		}

	} else {
		log.WithFields(log.Fields{
			"方法名": "GetShopInfo",
		}).Info("查询缓存成功")

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
	//清除缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "shopInfo", req.ShopId))

	var shop models.Shop
	result := s.db.Where("id = ?", req.ShopId).First(&shop)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
		ShopId:      req.ShopId,
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Picture:     req.Product.Picture,
		Price:       req.Product.Price,
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
			Id:          product.ID,
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

	return &pb.AddProductResp{ProductId: product.ID}, nil
}

// DeleteProduct 删除商品
func (s *ShopServiceImpl) DeleteProduct(ctx context.Context, req *pb.DeleteProductReq) (*pb.DeleteProductResp, error) {

	//清除缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(product, "item", req.ProductId))
	defer DelShopProductCache(ctx, rds, req.ShopId)

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

	//清除缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(product, "item", req.Product.Id))

	var product models.Product
	result := s.db.Where("id = ? AND shop_id = ?", req.Product.Id, req.ShopId).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &pb.UpdateProductResp{Res: false}, nil
		}
		return nil, result.Error
	}
	// 将 req.Product.Categories 从 []string 转换为 string
	categoriesStr := strings.Join(req.Product.Categories, ",")
	product.Name = req.Product.Name
	product.Description = req.Product.Description
	product.Picture = req.Product.Picture
	product.Price = req.Product.Price
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

	if req.Page < 1 || req.PageSize < 1 {
		return nil, BadPageOrPageSize
	}

	offset := (req.Page - 1) * req.PageSize

	log.Printf("Received GetProductList request for shop ID: %d, page: %d, page size: %d", req.ShopId, req.Page, req.PageSize)

	products := make([]models.Product, 0)

	//使用Md5作为条件,支持复杂扩容
	key := util.TakeKey(product, req.ShopId, util.Md5Hash(util.TakeKey(req.Page, req.PageSize)))

	result, _ := rds.Get(ctx, key).Result()

	list := make([]uint64, 0)
	fail := make([]uint64, 0)
	productMap := make(map[uint64]models.Product)

	//查询id数组
	ok := json.Unmarshal([]byte(result), &list)

	//按照list查询详情缓存
	for _, v := range list {

		t := models.Product{}
		key := util.TakeKey(product, "item", v)
		jsonData, err := rds.Get(ctx, key).Result()
		if err != nil || json.Unmarshal([]byte(jsonData), &t) != nil {
			fail = append(fail, v)
			continue
		}
		productMap[v] = t
	}

	//计算缓存失效比率
	rate := 100
	if len(list) != 0 {
		rate = len(fail) * 100 / len(list)
	}

	//查询缓存失败
	if ok != nil || len(list) == 0 || rate > 30 {

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

		//构建订单ID数组
		idList := make([]uint64, 0, len(products))
		for _, v := range products {
			idList = append(idList, v.ID)
		}

		//序列化并存储ID列表缓存
		jsonData, err := json.Marshal(&idList)

		if err != nil {
			util.LogError("序列化商品数组错误", "GetProductList", "请求hash为"+key, err)
		} else {
			// 设置随机过期时间，3~6 分钟
			if setErr := rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("存入商品缓存失败", "GetProductList", "", setErr)
			}
		}

		//订单详情分级存储
		for _, v := range products {
			key := util.TakeKey(product, "item", v.ID)
			jsonData, err := json.Marshal(&v)
			if err != nil {
				util.LogError("序列化商品错误", "GetProductList", "商品id为"+strconv.Itoa(int(v.ID)), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存商品失败", "GetProductList", "商品id为"+strconv.Itoa(int(v.ID)), err)
				}
			}
		}

	} else {

		//缓存失效比例较低,逐条查询并放入缓存
		if len(fail) > 0 {
			var missedProducts []models.Product
			if err := s.db.Where("id IN ?", fail).Find(&missedProducts).Error; err != nil {
				util.LogError("查询商品错误", "GetProductList", "", err)
				return nil, FailFetchProductList
			}
			for _, o := range missedProducts {
				productMap[o.ID] = o
				// 放入缓存
				key := util.TakeKey(product, "item", o.ID)
				jsonData, err := json.Marshal(&o)
				if err != nil {
					util.LogError("序列化商品错误", "GetProductList", "商品id为"+strconv.Itoa(int(o.ID)), err)
				} else {
					if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
						util.LogError("缓存商品失败", "GetProductList", "商品id为"+strconv.Itoa(int(o.ID)), err)
					}
				}
			}
		}

		// 最终按list顺序组装orders
		products = make([]models.Product, 0, len(list))
		for _, id := range list {
			if v, ok := productMap[id]; ok {
				products = append(products, v)
			}
		}

		log.WithFields(log.Fields{
			"方法名": "GetProductList",
		}).Info("查询缓存成功")
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
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Categories:  category,
		}
	}

	log.Printf("Successfully fetched %d products for shop ID %d", len(pbProducts), req.ShopId)
	return &pb.GetProductListResp{Products: pbProducts}, nil
}
