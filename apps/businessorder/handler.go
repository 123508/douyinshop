package main

import (
	"context"
	"encoding/json"
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
	"time"
)

// OrderBusinessServiceImpl implements the last service interface defined in the IDL.
type OrderBusinessServiceImpl struct{}

const (
	serviceName = "businessOrder"
)

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	return DB
}

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	return rds
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

//凡是只涉及查询功能的不应该上锁

// GetOrderList implements the OrderBusinessServiceImpl interface.
// 获取订单
func (s *OrderBusinessServiceImpl) GetOrderList(ctx context.Context, req *businessOrder.GetOrderListReq) (resp *businessOrder.GetOrderListResp, err error) {

	key := util.TakeKey(serviceName, "order", req.ShopId)

	result, err := rds.HGetAll(ctx, key).Result()

	if err != nil {
		util.LogError("查询缓存失败", "GetOrderList", "查询order缓存", err)
	}

	//查询订单
	orders := make([]models.Order, 0)

	for _, v := range result {
		o := models.Order{}
		err = json.Unmarshal([]byte(v), &o)
		if err != nil {
			util.LogError("反序列化数据失败", "GetOrderList", "反序列化order", err)
			orders = make([]models.Order, 0)
			break
		}
		orders = append(orders, o)
	}

	if err != nil { //查询缓存失败
		offset := (req.Page - 1) * req.PageSize

		//查询异常,返回错误
		if err = DB.Where("shop_id = ?", req.ShopId).Offset(int(offset)).Find(&orders).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderError
		}

		cache := map[string]string{}

		for _, v := range orders {
			jsonData, _ := json.Marshal(&v)
			cache[util.TakeKey(v.ID)] = string(jsonData)
		}

		//存储哈希
		err := rds.HMSet(ctx, key, cache).Err()

		if err != nil {
			util.LogError("缓存数据失败", "GetOrderList", "将order加入缓存", err)
		}

		// 设置随机过期时间，30~45 分钟
		rds.Expire(ctx, key, time.Duration(rand.Intn(15)+30)*time.Minute)

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

	orderKey := util.TakeKey(serviceName, "order", req.ShopId)

	//查询order缓存
	jsonData, err := rds.HGet(ctx, orderKey, util.TakeKey(req.OrderId)).Result()

	if err != nil {
		util.LogError("查询order缓存失败", "Detail", "查询order缓存", err)
	}

	err = json.Unmarshal([]byte(jsonData), &order)

	if err != nil { //查询缓存失败,直接去查数据库

		// 如果没有找到对应订单，返回错误信息
		if err = DB.Where("id = ? and shop_id = ?", req.OrderId, req.ShopId).First(&order).Error; err != nil {
			util.LogError("没有找到对应订单", "Detail", "查询order数据库", err)
			return nil, SearchOrderError
		}
	}

	// 查询订单详情（List）
	var orderDetails []models.OrderDetail

	detailKey := util.TakeKey(serviceName, "orderDetail", req.OrderId)

	//查询orderDetails缓存
	result, err := rds.HGetAll(ctx, detailKey).Result()

	if err != nil {
		util.LogError("查询orderDetail缓存失败", "Detail", "查询orderDetail缓存", err)
	}

	for _, v := range result {

		t := models.OrderDetail{}

		err = json.Unmarshal([]byte(v), &t)

		if err != nil {
			util.LogError("反序列化失败", "Detail", "查询orderDetail", err)
			orderDetails = make([]models.OrderDetail, 0)
			break
		}

		orderDetails = append(orderDetails, t)
	}

	if err != nil {

		// 如果查询订单详情失败，返回错误
		if err = DB.Where("order_id = ?", req.OrderId).Find(&orderDetails).Error; err != nil {
			util.LogError("查询订单详情失败", "Detail", "查询订单详情", err)
			return nil, SearchOrderDetailsError
		}

		cache := map[string]string{}

		for _, v := range orderDetails {
			jsonData, err := json.Marshal(&v)
			if err != nil {
				util.LogError("序列化数据失败", "Detail", "序列化orderDetail", err)
			}
			cache[util.TakeKey(v.ID)] = string(jsonData)
		}

		//存储哈希
		err := rds.HMSet(ctx, detailKey, cache).Err()

		if err != nil {
			util.LogError("缓存数据失败", "Detail", "将orderDetail加入缓存", err)
		}

		// 设置随机过期时间，30~45 分钟
		rds.Expire(ctx, detailKey, time.Duration(rand.Intn(15)+30)*time.Minute)
	}

	//查询订单日志详情
	var orderLogs []models.OrderStatusLog

	logKey := util.TakeKey(serviceName, "orderLog", req.OrderId)

	res, err := rds.HGetAll(ctx, logKey).Result()

	if err != nil {
		util.LogError("查询orderLog缓存失败", "Detail", "查询order缓存", err)
	}

	for _, v := range res {
		l := models.OrderStatusLog{}
		err = json.Unmarshal([]byte(v), &l)
		if err != nil {
			util.LogError("反序列化失败", "Detail", "查询orderLog", err)
			orderLogs = make([]models.OrderStatusLog, 0)
			break
		}
		orderLogs = append(orderLogs, l)
	}

	if err != nil { //查询缓存失败
		//查询失败报错
		if err = DB.Where("order_id = ?", req.OrderId).Find(&orderLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogsError
		}

		cache := map[string]string{}

		for _, v := range orderLogs {
			jsonData, err := json.Marshal(&v)

			if err != nil {
				util.LogError("序列化数据失败", "Detail", "序列化orderLog", err)
			}

			cache[util.TakeKey(v.ID)] = string(jsonData)
		}

		err := rds.HMSet(ctx, detailKey, cache).Err()

		if err != nil {
			util.LogError("缓存数据失败", "Detail", "将orderLog加入缓存", err)
		}

		rds.Expire(ctx, detailKey, time.Duration(rand.Intn(15)+30)*time.Minute)
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
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.ShopId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.OrderId))

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
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.ShopId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.OrderId))

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
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.ShopId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.OrderId))

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
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.ShopId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.OrderId))

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
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.ShopId))
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, "order", req.OrderId))

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
