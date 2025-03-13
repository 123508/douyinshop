package main

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// OrderUserServiceImpl implements the last service interface defined in the IDL.
type OrderUserServiceImpl struct{}

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
		fmt.Println(err)
		return nil, SearchOrderError
	}
	return &order, nil
}

var SubmitOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单提交失败"}

var SearchOrderError = &errorno.BasicMessageError{Code: 404, Message: "查询订单错误"}

var SearchOrderLogError = &errorno.BasicMessageError{Code: 404, Message: "查询订单日志错误"}

var CancelOrderError = &errorno.BasicMessageError{Code: 500, Message: "取消订单失败"}

var UnableChangeStatusError = &errorno.BasicMessageError{Code: 400, Message: "更新状态失败,该状态不允许被更新"}

var NotPayReminderError = &errorno.BasicMessageError{Code: 400, Message: "没有支付订单,无法提醒发货"}

var DeliveredReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家已经发货,无需提醒"}

var CancelReminderError = &errorno.BasicMessageError{Code: 400, Message: "订单已经取消,无法提醒发货"}

var RefundingReminderError = &errorno.BasicMessageError{Code: 400, Message: "退款中,无法提醒发货"}

var RefundedReminderError = &errorno.BasicMessageError{Code: 400, Message: "已经退款,无法提醒发货"}

var RejectionReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家拒绝发货,无法提醒"}

var StatusError = &errorno.BasicMessageError{Code: 400, Message: "不允许的行为"}

var CompletionOrderError = &errorno.BasicMessageError{Code: 400, Message: "无法确认收货"}

// Submit implements the OrderUserServiceImpl interface.
// 用户提交订单
func (s *OrderUserServiceImpl) Submit(ctx context.Context, req *userOrder.OrderSubmitReq) (resp *userOrder.OrderSubmitResp, err error) {

	// 创建订单对象
	var order models.Order

	err = DB.Transaction(func(tx *gorm.DB) error {

		order.Number = uuid.New().String() // 使用 UUID 生成唯一的订单号
		order.UserId = req.UserId
		order.AddressBookId = int(req.AddressBookId)
		order.PayMethod = int(req.PayMethod)
		order.Remark = req.Remark
		order.ShopId = req.Order.ShopId

		var total float64

		for _, detail := range req.Order.List {
			total += float64(detail.Amount) * float64(detail.Number)
		}

		order.Amount = total

		if err = DB.Create(&order).Error; err != nil {
			return err
		}

		var orderDetails []models.OrderDetail
		for _, detail := range req.Order.List {

			orderDetail := models.OrderDetail{
				Name:      detail.Name,
				Image:     detail.Image,
				OrderId:   int(order.ID),
				ProductId: int(detail.ProductId),
				Number:    int(detail.Number),
				Amount:    float64(detail.Amount),
			}

			// 将订单详情添加到列表中
			orderDetails = append(orderDetails, orderDetail)
		}

		//如果创建失败就返回错误
		if err = DB.Create(&orderDetails).Error; err != nil {
			return err
		}

		current := time.Now()

		task := time.Now().Add(15 * time.Minute)

		orderStatusLog := models.OrderStatusLog{
			OrderId:     int(order.ID),
			Status:      0, // 初始状态为待付款
			StartTime:   &current,
			EndTime:     &task,
			Description: "订单创建，待付款",
		}

		// 保存状态日志
		if err = DB.Create(&orderStatusLog).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, SubmitOrderError
	}

	// 构建返回对象
	resp = &userOrder.OrderSubmitResp{
		OrderId:     uint32(order.ID),
		Number:      order.Number,
		OrderAmount: float32(order.Amount),
	}

	return resp, nil
}

// History implements the OrderUserServiceImpl interface.
// 查询用户的历史订单
func (s *OrderUserServiceImpl) History(ctx context.Context, req *userOrder.HistoryReq) (resp *userOrder.HistoryResp, err error) {

	// 检查分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	offset := (req.Page - 1) * req.PageSize

	// 查询用户订单数据
	var orders []models.Order
	err = DB.Where("user_id = ?", req.UserId).
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find(&orders).Error

	if err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	// 构建订单响应数据
	orderList := make([]*order_common.OrderResp, len(orders))
	for i, order := range orders {
		// 根据订单号查询订单详情
		var orderDetails []models.OrderDetail

		if err = DB.Where("order_id = ?", order.ID).Find(&orderDetails).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderError
		}

		var orderLog models.OrderStatusLog

		if err = DB.Where("order_id = ?", order.ID).Last(&orderLog).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogError
		}

		// 构建订单响应
		orderResp := &order_common.OrderResp{
			Order: &order_common.Order{
				ID:          uint32(order.ID),
				Number:      order.Number,
				PayStatus:   int32(order.PayStatus),
				Amount:      float32(order.Amount),
				FinalStatus: uint32(orderLog.Status),
			},
		}
		orderList[i] = orderResp
	}

	// 返回响应
	return &userOrder.HistoryResp{
		List:     orderList,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    uint32(len(orders)),
	}, nil
}

// Detail implements the OrderUserServiceImpl interface.
// 查询订单的详细信息
func (s *OrderUserServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {

	order, err := GetOrderInfo(req.OrderId)

	if err != nil {
		return nil, err
	}

	// 查询订单详情（List）
	var orderDetails []models.OrderDetail

	// 如果查询订单详情失败，返回错误
	if err = DB.Where("order_id = ?", req.OrderId).Find(&orderDetails).Error; err != nil {
		fmt.Println(err)
		return nil, SearchOrderError
	}

	var orderLogs []models.OrderStatusLog

	if err = DB.Where("order_id = ?", order.ID).Find(&orderLogs).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogError
	}

	var statusList []*order_common.Status

	for _, k := range orderLogs {
		statusList = append(statusList, &order_common.Status{
			StartTime: k.StartTime.String(),
			Status:    uint32(k.Status),
			EndTime:   k.EndTime.String(),
		})
	}

	// 处理类型
	orderCommon := &order_common.Order{
		Number:        order.Number,
		UserId:        order.UserId,
		PayMethod:     int32(order.PayMethod),
		PayStatus:     int32(order.PayStatus),
		AddressBookId: uint64(order.AddressBookId),
		Amount:        float32(order.Amount),
		Remark:        order.Remark,
		Phone:         order.Phone,
		Address:       order.Address,
		Username:      order.UserName,
		Consignee:     order.Consignee,
		Status:        statusList,
	}

	var orderDetail []*order_common.OrderDetail

	for _, k := range orderDetails {
		orderDetail = append(orderDetail, &order_common.OrderDetail{
			Name:      k.Name,
			Image:     k.Image,
			OrderId:   uint32(k.OrderId),
			ProductId: uint32(k.ProductId),
			Number:    uint32(k.Number),
			Amount:    float32(k.Amount),
		})
	}

	// 构建并返回响应数据
	resp = &order_common.OrderResp{
		Order:        orderCommon,
		OrderDetails: orderDetail,
	}

	return resp, nil
}

// Cancel implements the OrderUserServiceImpl interface.
// 取消订单
// 步骤:查询要取消的订单->查询不到，返回错误，否则继续->判断该订单的status是否<5,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	//查询要取消的订单
	_, err = GetOrderInfo(req.OrderId)

	//查询要取消的订单失败
	if err != nil {
		return nil, err
	}

	//查询订单状态
	var status models.OrderStatusLog

	//查询订单日志异常
	if err = DB.Where("order_id = ?", req.OrderId).Last(&status).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogError
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	//订单状态错误
	if status.Status > 5 {
		return nil, UnableChangeStatusError
	}

	//决定下一个状态
	var Status uint

	var Description string

	if status.Status == 0 {
		Status = 6
		Description = "已取消"
	} else {
		Status = 7
		Description = "退款中"
	}

	//创建新状态
	newStatus := models.OrderStatusLog{
		StartTime:   &currentTime,
		Status:      Status,
		EndTime:     nil,
		Description: Description,
	}

	// 更新订单详情的状态为“已取消”并记录到 OrderStatusLog
	err = DB.Transaction(func(tx *gorm.DB) error {

		//将原有状态的结束时间修改
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", status.ID).Update("end_time", currentTime).Error; err != nil {
			return err
		}

		//插入新的状态
		if err = DB.Create(&newStatus).Error; err != nil {
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

// Reminder implements the OrderUserServiceImpl interface.
// 提醒商家发货
// 步骤:查询订单是否属于用户->判断status==2,不是就返回错误,是继续(增加健壮性)->通过订单查询商家id,查不到就报错,否则继续->发送信息
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {

	// 查询该用户的订单信息
	var order models.Order

	//查询订单失败,返回异常
	if err = DB.Where("id = ? AND user_id = ?", req.OrderId, req.UserId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderError
	}

	//查询订单状态信息
	var status models.OrderStatusLog

	//查询失败,返回异常
	if err = DB.Where("order_id = ?", req.OrderId).Last(&status).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogError
	}

	//如果状态不是商家已接单就返回错误

	switch status.Status {
	case 0:
		return nil, NotPayReminderError
	case 1, 2: //只允许这两个状态通行
	case 3, 4, 5:
		return nil, DeliveredReminderError
	case 6:
		return nil, CancelReminderError
	case 7:
		return nil, RefundingReminderError
	case 8:
		return nil, RefundedReminderError
	case 9:
		return nil, RejectionReminderError
	default:
		return nil, StatusError
	}

	//提醒商家发货逻辑
	err = util.SendMessage("order.direct", "message", "orderId:"+strconv.Itoa(int(req.OrderId))+",shopId:"+strconv.Itoa(int(order.ShopId)), 1)

	if err != nil {
		return nil, err
	}

	return &order_common.Empty{}, nil
}

// Complete implements the OrderUserServiceImpl interface.
// 确认收货
// 订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 取消退款(直接回到上一步即可)
// 步骤:查询要完成的订单->查询不到，返回错误，否则继续->判断该订单的status是否为4,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Complete(ctx context.Context, req *userOrder.CompleteReq) (resp *order_common.Empty, err error) {
	// 查询订单信息
	_, err = GetOrderInfo(req.OrderId)

	//查询不到订单信息,返回异常
	if err != nil {
		return nil, err
	}

	//查询订单状态信息
	var status models.OrderStatusLog

	//查询订单日志失败,返回异常
	if err = DB.Where("order_id = ?", req.OrderId).Last(&status).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogError
	}

	//如果订单状态不为待派送就无法收货
	if status.Status != 4 && status.Status != 3 {
		return nil, UnableChangeStatusError
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	//创建完成状态
	newStatus := models.OrderStatusLog{
		StartTime:   &currentTime,
		Status:      5,
		EndTime:     nil,
		Description: "已完成",
	}
	err = DB.Transaction(func(tx *gorm.DB) error {

		//将原有状态的结束时间修改
		if err = DB.Model(&models.OrderStatusLog{}).Where("id = ?", status.ID).Update("end_time", currentTime).Error; err != nil {
			return err
		}

		//插入新的状态
		if err = DB.Create(&newStatus).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, CompletionOrderError
	}

	return &order_common.Empty{}, nil
}
