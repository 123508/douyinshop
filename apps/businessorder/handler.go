package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

// OrderBusinessServiceImpl implements the last service interface defined in the IDL.
type OrderBusinessServiceImpl struct{}

const (
	goods       = "goods"
	serviceName = "businessOrder"
)

//订单列表的缓存方案 : id列表+详情分级缓存

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", "", err)
	}
	return DB
}

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

func cleanCache(ctx context.Context, orderId uint32) {
	//清除缓存
	util.CleanCache(rds, ctx, util.TakeKey(goods, "order", orderId))
	util.CleanCache(rds, ctx, util.TakeKey(goods, "orderDetail", orderId))
	util.CleanCache(rds, ctx, util.TakeKey(goods, "orderLog", orderId))
}

// GetOrderInfo 查询订单信息
func GetOrderInfo(orderId uint32) (*models.Order, error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	// 如果没有找到对应订单，返回错误信息
	if err := DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}
	return &order, nil
}

// 核心方法，用于更新订单的状态
func (s *OrderBusinessServiceImpl) updateOrderStatus(orderId uint32, status uint32) (err error) {

	order, err := GetOrderInfo(orderId)

	if err != nil {
		return err
	}

	var orderLog models.OrderStatusLog

	if err = DB.Model(&models.OrderStatusLog{}).Where("order = ? and version = ?", orderId, order.FinalVersion).Last(&orderLog).Error; err != nil {
		log.Println(err)
		return SearchOrderLogsError
	}

	// 订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 取消退款(直接回到上一步即可)
	if status == 3 {
		if order.FinalStatus != 2 {
			return
		}
	}

	current := time.Now()

	newOrderLog := models.OrderStatusLog{
		OrderId:   orderId,
		Status:    status,
		StartTime: &current,
		EndTime:   nil,
		Version:   orderLog.Version + 1,
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		//更新旧订单状态
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", orderLog.ID).Update("end_time", &current).Error; err != nil {
			return err
		}

		//插入新订单状态
		if err = DB.Create(&newOrderLog).Error; err != nil {
			return err
		}

		//更新最终状态
		if err = DB.Model(&models.Order{}).Where("id = ?", orderId).Update("final_status", newOrderLog.Status).Update("final_version", newOrderLog.Version).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

var NoOrderIdError = &errorno.BasicMessageError{Code: 404, Message: "没有用户id"}

var SearchOrderError = &errorno.BasicMessageError{Code: 500, Message: "查询订单错误"}

var SearchOrderDetailsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单详细信息错误"}

var SearchOrderLogsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单日志信息错误"}

var CancelOrderError = &errorno.BasicMessageError{Code: 500, Message: "取消订单错误"}

var DeliveryOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单发货失败"}

var ReceiveOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单变为待收货失败"}

var RejectionOrderError = &errorno.BasicMessageError{Code: 500, Message: "拒单失败"}

var ConfirmOrderError = &errorno.BasicMessageError{Code: 500, Message: "确认订单失败"}

var UnableRejectionOrderError = &errorno.BasicMessageError{Code: 500, Message: "当前状态不允许拒单,请注意"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}

//凡是只涉及查询功能的不应该上锁

// GetOrderList implements the OrderBusinessServiceImpl interface.
// 获取订单
func (s *OrderBusinessServiceImpl) GetOrderList(ctx context.Context, req *businessOrder.GetOrderListReq) (resp *businessOrder.GetOrderListResp, err error) {

	if req.Page < 1 || req.PageSize < 1 {
		return nil, BadPageOrPageSize
	}

	//查询订单
	orders := make([]models.Order, 0)

	//使用Md5作为条件,支持复杂扩容
	key := util.Md5Hash(util.TakeKey(serviceName, "list", req.ShopId, req.Page, req.PageSize))

	result, err := rds.Get(ctx, key).Result()

	if err != nil {
		util.LogError("查询缓存失败", "GetOrderList", "查询order缓存", err)
	}

	list := make([]uint32, 0)
	fail := make([]uint32, 0)
	orderMap := make(map[uint32]models.Order)

	//查询id数组
	ok := json.Unmarshal([]byte(result), &list)

	//查询数组为空的时候直接进行返回处理
	if len(list) == 0 {
		return &businessOrder.GetOrderListResp{
			List: make([]*order_common.OrderResp, 0),
		}, nil
	}

	for _, v := range list {

		t := models.Order{}
		key := util.TakeKey(goods, "order", v)
		jsonData, err := rds.Get(ctx, key).Result()
		if err != nil || json.Unmarshal([]byte(jsonData), &t) != nil {
			fail = append(fail, v)
			continue
		}
		orderMap[v] = t
	}

	//计算缓存失效比率
	rate := len(fail) * 100 / len(list)

	//查询缓存失败或者缓存失效比率过大
	if ok != nil || len(list) == 0 || rate > 30 {
		orders = make([]models.Order, 0)
		offset := (req.Page - 1) * req.PageSize

		//查询异常,返回错误
		if err = DB.Where("shop_id = ?", req.ShopId).Offset(int(offset)).Limit(int(req.PageSize)).Find(&orders).Error; err != nil {
			util.LogError("查询order错误", "GetOrderList", "", err)
			return nil, SearchOrderError
		}

		//构建商品ID数组
		idList := make([]uint32, 0, len(orders))
		for _, v := range orders {
			idList = append(idList, v.ID)
		}

		jsonData, err := json.Marshal(&idList)

		if err != nil {
			util.LogError("序列化商品数组错误", "GetOrderList", "请求hash为"+key, err)
		} else {
			// 设置随机过期时间，3~6 分钟
			if setErr := rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("存入order缓存失败", "GetOrderList", "", setErr)
			}
		}

		//商品详情分级存储
		for _, v := range orders {
			key := util.TakeKey(goods, "order", v.ID)
			jsonData, err := json.Marshal(&v)
			if err != nil {
				util.LogError("序列化商品错误", "GetOrderList", "order的id为"+strconv.Itoa(int(v.ID)), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存商品失败", "GetOrderList", "order的id为"+strconv.Itoa(int(v.ID)), err)
				}
			}
		}

	} else {

		//缓存失效比例较低,逐条查询并放入缓存
		if len(fail) > 1 {
			var missedOrders []models.Order
			if err := DB.Where("id IN ?", fail).Find(&missedOrders).Error; err != nil {
				util.LogError("查询order错误", "GetOrderList", "", err)
				return nil, SearchOrderError
			}
			for _, o := range missedOrders {
				orderMap[o.ID] = o
				//放入缓存
				key := util.TakeKey(goods, "order", o.ID)

				jsonData, err := json.Marshal(&o)
				if err != nil {
					util.LogError("序列化商品错误", "GetOrderList", "商品id为"+strconv.Itoa(int(o.ID)), err)
				} else {
					if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
						util.LogError("缓存商品失败", "GetOrderList", "商品id为"+strconv.Itoa(int(o.ID)), err)
					}
				}
			}
		} else if len(fail) == 1 {
			t := models.Order{}

			//查询异常,返回错误
			if err := DB.Where("id = ?", fail[0]).First(&t).Error; err != nil {
				util.LogError("查询order错误", "GetOrderList", "", err)
				return nil, SearchOrderError
			}

			orderMap[fail[0]] = t

			//放入缓存
			key := util.TakeKey(goods, "order", fail[0])

			jsonData, err := json.Marshal(&t)
			if err != nil {
				util.LogError("序列化商品错误", "GetOrderList", "商品id为"+strconv.Itoa(int(fail[0])), err)
			} else {
				if err = rds.Set(ctx, key, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); err != nil {
					util.LogError("缓存商品失败", "GetOrderList", "商品id为"+strconv.Itoa(int(fail[0])), err)
				}
			}
		}

		// 最终按list顺序组装orders
		orders = make([]models.Order, 0, len(list))
		for _, id := range list {
			if v, ok := orderMap[id]; ok {
				orders = append(orders, v)
			}
		}

		log.WithFields(log.Fields{
			"方法名": "GetOrderList",
		}).Info("查询order缓存成功")
	}

	// 封装订单响应列表
	var orderRespList []*order_common.OrderResp
	for _, order := range orders {
		// 将 models.Order 转换为 order_common.Order
		orderCommon := &order_common.Order{
			Number:      order.Number,
			PayMethod:   order.PayMethod,
			Amount:      order.Amount,
			Phone:       order.Phone,
			FinalStatus: order.FinalStatus,
		}

		// 封装每个订单的响应
		orderResp := &order_common.OrderResp{
			Order: orderCommon,
		}

		orderRespList = append(orderRespList, orderResp)
	}

	// 返回订单列表响应
	return &businessOrder.GetOrderListResp{
		List: orderRespList,
	}, nil
}

// Detail implements the OrderBusinessServiceImpl interface.
// 订单详情
func (s *OrderBusinessServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {
	// 查询订单信息（Order）基本信息
	var order models.Order

	orderKey := util.TakeKey(goods, "order", req.OrderId)

	//查询单条order缓存
	jsonData, err := rds.Get(ctx, orderKey).Result()

	if errors.Is(err, redis.Nil) {
		log.Println("order缓存未命中")
	} else if err != nil {
		util.LogError("查询order缓存失败", "Detail", "查询order缓存", err)
	}

	//查询缓存失败,直接去查数据库
	if err != nil || json.Unmarshal([]byte(jsonData), &order) != nil {

		// 如果没有找到对应订单，返回错误信息
		if err = DB.Where("id = ? and shop_id = ?", req.OrderId, req.ShopId).First(&order).Error; err != nil {
			util.LogError("没有找到对应订单", "Detail", "查询order数据库", err)
			return nil, SearchOrderError
		}

		jsonData, err := json.Marshal(&order)

		if err != nil {
			util.LogError("序列化商品数据错误", "Detail", "", err)
		} else {
			//将数据加入缓存
			if setErr := rds.Set(ctx, orderKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "Detail", "将orderDetail加入缓存", setErr)
			}
		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "Detail",
		}).Info("查询order缓存成功")
	}

	// 查询订单详情（List）
	orderDetails := make([]models.OrderDetail, 0)

	detailKey := util.TakeKey(goods, "orderDetail", req.OrderId)

	//查询orderDetails缓存
	result, err := rds.Get(ctx, detailKey).Result()

	if err != nil {
		util.LogError("查询orderDetail缓存失败", "Detail", "查询orderDetail缓存", err)
	}

	//查询orderDetail缓存失败
	if err != nil || json.Unmarshal([]byte(result), &orderDetails) != nil {

		// 如果查询订单详情失败，返回错误
		if dbErr := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails).Error; dbErr != nil {
			util.LogError("查询订单详情失败", "Detail", "查询订单详情", dbErr)
			return nil, SearchOrderDetailsError
		}

		jsonData, err := json.Marshal(&orderDetails)

		if err != nil {
			util.LogError("序列化数据失败", "Detail", "序列化orderDetail", err)
		} else {
			if setErr := rds.Set(ctx, detailKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "Detail", "将orderDetail加入缓存", setErr)
			}
		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "Detail",
		}).Info("查询orderDetail缓存成功")
	}

	//查询订单日志详情
	orderLogs := make([]models.OrderStatusLog, 0)

	logKey := util.TakeKey(goods, "orderLog", req.OrderId)

	res, err := rds.Get(ctx, logKey).Result()

	if err != nil {
		util.LogError("查询orderLog缓存失败", "Detail", "查询order缓存", err)
	}

	//查询缓存失败
	if err != nil || json.Unmarshal([]byte(res), &orderLogs) != nil {
		//查询失败报错
		if err = DB.Where("order_id = ?", req.OrderId).Find(&orderLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogsError
		}

		jsonData, err := json.Marshal(&orderLogs)

		if err != nil {
			util.LogError("序列化数据失败", "Detail", "序列化orderLog", err)
		} else {
			if setErr := rds.Set(ctx, logKey, jsonData, time.Duration(rand.Intn(3)+3)*time.Minute).Err(); setErr != nil {
				util.LogError("缓存数据失败", "Detail", "将orderLog加入缓存", setErr)
			}

		}
	} else {
		log.WithFields(log.Fields{
			"方法名": "Detail",
		}).Info("查询orderLog缓存成功")
	}

	//将切片中的OrderStatusLog转换为order_common.Status
	var newStatus []*order_common.Status
	for _, orderLog := range orderLogs {

		newStatus = append(newStatus, &order_common.Status{
			StartTime:   orderLog.StartTime.String(),
			EndTime:     orderLog.EndTime.String(),
			Status:      orderLog.Status,
			Description: orderLog.Description,
		})
	}

	// 封装订单
	orderCommon := &order_common.Order{
		Number:        order.Number,
		UserId:        order.UserId,
		PayMethod:     order.PayMethod,
		AddressBookId: order.AddressBookId,
		Amount:        order.Amount,
		Remark:        order.Remark,
		Phone:         order.Phone,
		Address:       order.Address,
		Username:      order.UserName,
		Consignee:     order.Consignee,
		Status:        newStatus,
	}

	//封装订单明细
	var Details []*order_common.OrderDetail
	for _, orderDetail := range orderDetails {
		Details = append(Details, &order_common.OrderDetail{
			Name:      orderDetail.Name,
			Image:     orderDetail.Image,
			OrderId:   orderDetail.OrderId,
			ProductId: orderDetail.ProductId,
			Number:    orderDetail.Number,
			Amount:    orderDetail.Amount,
		})
	}

	// 构建并返回响应数据
	return &order_common.OrderResp{
		Order:        orderCommon,
		OrderDetails: Details,
	}, nil
}

// Confirm implements the OrderBusinessServiceImpl interface.
// 确认订单
func (s *OrderBusinessServiceImpl) Confirm(ctx context.Context, req *businessOrder.ConfirmReq) (resp *order_common.Empty, err error) {

	//清除缓存
	defer cleanCache(ctx, req.OrderId)

	if err = s.updateOrderStatus(req.OrderId, 2); err != nil {
		util.LogError("修改订单失败", "Confirm", "", err)
		return nil, ConfirmOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Delivery implements the OrderBusinessServiceImpl interface.
// 商家发货
func (s *OrderBusinessServiceImpl) Delivery(ctx context.Context, req *businessOrder.DeliveryReq) (resp *order_common.Empty, err error) {

	//清除缓存
	defer cleanCache(ctx, req.OrderId)

	if err = s.updateOrderStatus(req.OrderId, 3); err != nil {
		util.LogError("修改订单失败", "Delivery", "", err)
		return nil, DeliveryOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Receive implements the OrderBusinessServiceImpl interface.
// 转变订单状态为待收货
func (s *OrderBusinessServiceImpl) Receive(ctx context.Context, req *businessOrder.ReceiveReq) (resp *order_common.Empty, err error) {

	//清除缓存
	defer cleanCache(ctx, req.OrderId)

	if err = s.updateOrderStatus(req.OrderId, 4); err != nil {
		log.Println(err)
		util.LogError("修改订单失败", "Receive", "", err)
		return nil, ReceiveOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Rejection implements the OrderBusinessServiceImpl interface.
// 商家拒绝订单
// 订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 取消退款(直接回到上一步即可)
func (s *OrderBusinessServiceImpl) Rejection(ctx context.Context, req *businessOrder.RejectionReq) (resp *order_common.Empty, err error) {

	//清除缓存
	defer cleanCache(ctx, req.OrderId)

	order, err := GetOrderInfo(req.OrderId)

	if err != nil {
		return nil, err
	}

	var orderLog models.OrderStatusLog

	if err = DB.Where("order_id = ? and version = ?", req.OrderId, order.FinalVersion).Last(&orderLog).Error; err != nil {
		util.LogError("查询数据库失败", "Rejection", "查询order订单", err)
		return nil, SearchOrderLogsError
	}

	if order.FinalStatus > 2 {
		return nil, UnableRejectionOrderError
	}

	var status uint32

	current := time.Now()

	if orderLog.Status == 2 {
		status = 7
	} else {
		status = 9
	}

	newOrderLog := models.OrderStatusLog{
		OrderId:     req.OrderId,
		Status:      status,
		StartTime:   &current,
		EndTime:     nil,
		Description: req.RejectionReason,
		Version:     orderLog.Version + 1,
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		//更新旧订单状态
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", orderLog.ID).Update("end_time", &current).Error; err != nil {
			return err
		}

		//插入新订单状态
		if err = DB.Create(&newOrderLog).Error; err != nil {
			return err
		}

		//更新最终状态
		if err = DB.Model(&models.Order{}).Where("id = ?", req.OrderId).Update("final_status", newOrderLog.Status).Update("final_version", newOrderLog.Version).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		util.LogError("更新订单失败", "Rejection", "更新orderLog", err)
		return nil, RejectionOrderError
	}

	return &order_common.Empty{}, nil
}

// Cancel implements the OrderBusinessServiceImpl interface.
// 商家取消订单(整体逻辑基本保持与用户端一样)
func (s *OrderBusinessServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	//清除缓存
	defer cleanCache(ctx, req.OrderId)

	order, err := GetOrderInfo(req.OrderId)

	if err != nil {
		return nil, err
	}

	var orderLog models.OrderStatusLog

	if err = DB.Model(&models.OrderStatusLog{}).Where("order = ? and version = ?", req.OrderId, order.FinalVersion).Last(&orderLog).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogsError
	}

	var status uint32

	if order.FinalStatus == 0 {
		status = 6
	} else {
		status = 7
	}

	current := time.Now()

	newOrderLog := models.OrderStatusLog{
		StartTime:   &current,
		Status:      status,
		EndTime:     nil,
		Description: req.CancelReason,
		Version:     orderLog.Version + 1,
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		//更新旧订单状态
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", orderLog.ID).Update("end_time", &current).Error; err != nil {
			return err
		}

		//插入新订单状态
		if err = DB.Create(&newOrderLog).Error; err != nil {
			return err
		}

		//更新最终状态
		if err = DB.Model(&models.Order{}).Where("id = ?", req.OrderId).Update("final_status", newOrderLog.Status).Update("final_version", newOrderLog.Version).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		util.LogError("更新订单失败", "Cancel", "", err)
		return nil, CancelOrderError
	}

	return &order_common.Empty{}, nil
}

// GetNotify implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) GetNotify(ctx context.Context, req *businessOrder.GetNotifyReq) (resp *businessOrder.GetNotifyResp, err error) {

	if req.OrderId == 0 {
		return nil, NoOrderIdError
	}

	n := &businessOrder.Notify{OrderId: req.OrderId, Content: "用户提醒发货"}

	return &businessOrder.GetNotifyResp{Notify: n}, nil
}
