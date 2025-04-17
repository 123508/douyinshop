package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"gorm.io/gorm"
	"log"
	"time"
)

// OrderBusinessServiceImpl implements the last service interface defined in the IDL.
type OrderBusinessServiceImpl struct{}

var DB = open()

func open() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	return DB
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

	//查询订单
	var orders []models.Order
	offset := (req.Page - 1) * req.PageSize

	//查询异常,返回错误
	if err = DB.Where("shop_id = ?", req.ShopId).Offset(int(offset)).Find(&orders).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
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

	// 如果没有找到对应订单，返回错误信息
	if err = DB.Where("id = ? and shop_id = ?", req.OrderId, req.ShopId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	// 查询订单详情（List）
	var orderDetails []models.OrderDetail

	// 如果查询订单详情失败，返回错误
	if err = DB.Where("order_id = ?", req.OrderId).Find(&orderDetails).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderDetailsError
	}

	//查询订单日志详情
	var orderLogs []models.OrderStatusLog

	//查询失败报错
	if err = DB.Where("order_id = ?", req.OrderId).Find(&orderLogs).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogsError
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

	if err = s.updateOrderStatus(req.OrderId, 2); err != nil {
		log.Println(err)
		return nil, ConfirmOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Delivery implements the OrderBusinessServiceImpl interface.
// 商家发货
func (s *OrderBusinessServiceImpl) Delivery(ctx context.Context, req *businessOrder.DeliveryReq) (resp *order_common.Empty, err error) {

	if err = s.updateOrderStatus(req.OrderId, 3); err != nil {
		log.Println(err)
		return nil, DeliveryOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Receive implements the OrderBusinessServiceImpl interface.
// 转变订单状态为待收货
func (s *OrderBusinessServiceImpl) Receive(ctx context.Context, req *businessOrder.ReceiveReq) (resp *order_common.Empty, err error) {

	if err = s.updateOrderStatus(req.OrderId, 4); err != nil {
		log.Println(err)
		return nil, ReceiveOrderError
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Rejection implements the OrderBusinessServiceImpl interface.
// 商家拒绝订单
// 订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 取消退款(直接回到上一步即可)
func (s *OrderBusinessServiceImpl) Rejection(ctx context.Context, req *businessOrder.RejectionReq) (resp *order_common.Empty, err error) {

	order, err := GetOrderInfo(req.OrderId)

	if err != nil {
		return nil, err
	}

	var orderLog models.OrderStatusLog

	if err = DB.Where("order_id = ? and version = ?", req.OrderId, order.FinalVersion).Last(&orderLog).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogsError
	}

	if order.FinalStatus > 2 {
		log.Println("当前状态无法拒单,请注意")
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
		log.Println(err)
		return nil, RejectionOrderError
	}

	return &order_common.Empty{}, nil
}

// Cancel implements the OrderBusinessServiceImpl interface.
// 商家取消订单(整体逻辑基本保持与用户端一样)
func (s *OrderBusinessServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

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
		log.Println(err)
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
