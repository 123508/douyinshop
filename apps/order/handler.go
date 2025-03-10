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

//订单状态 0待付款 1待接单 2已接单 3运输中 4待收货 5已完成 6已取消 7退款中 8已退款 9商家拒单 10取消退款(然后继续回到上一步)

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

var SaveLogError = &errorno.BasicMessageError{Code: 500, Message: "保存日志失败"}

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
		order.Amount = float64(req.Amount)

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

		if err = DB.Create(&orderDetails).Error; err != nil {
			return err
		}

		for _, detail := range orderDetails {
			orderStatusLog := models.OrderStatusLog{
				OrderDetailId: int(detail.ID),
				Status:        0, // 初始状态为待付款
				StartTime:     time.Now(),
				EndTime:       time.Now().Add(15 * time.Minute),
				Description:   "订单创建，待付款",
			}

			// 保存状态日志
			if err = DB.Create(&orderStatusLog).Error; err != nil {
				return err
			}
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

		var orderLogs models.OrderStatusLog

		if err = DB.Where("order_id = ?", order.ID).Last(&orderLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogError
		}

		// 构建订单响应
		orderResp := &order_common.OrderResp{
			Order: &order_common.Order{
				ID:        uint32(order.ID),
				Number:    order.Number,
				PayStatus: int32(order.PayStatus),
				Amount:    float32(order.Amount),
				Status:    int32(orderLogs.Status),
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

	var orderLogs models.OrderStatusLog

	if err = DB.Where("order_id = ?", order.ID).Last(&orderLogs).Error; err != nil {
		log.Println(err)
		return nil, SearchOrderLogError
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
		Status:        int32(orderLogs.Status),
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
// TODO 需要修改
// 步骤:查询要取消的订单->查询不到，返回错误，否则继续->判断该订单的status是否<5,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	_, err = GetOrderInfo(req.OrderId)

	if err != nil {
		return nil, err
	}

	// 查询订单详情（OrderDetails）
	var orderDetails []models.OrderDetail

	if err = DB.Where("order_id = ?", req.OrderId).Find(&orderDetails).Error; err != nil {
		return nil, err
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	// 更新订单详情的状态为“已取消”并记录到 OrderStatusLog
	for _, detail := range orderDetails {

		if err = DB.Save(&detail).Error; err != nil {
			return nil, err
		}

		// 查找并更新与当前 OrderDetailId 匹配的日志的 EndTime
		var statusLogs []models.OrderStatusLog
		statusLogResult := DB.Where("order_detail_id = ? AND start_time != end_time", detail.ID).Find(&statusLogs)

		if statusLogResult.Error != nil {
			return nil, statusLogResult.Error
		}

		// 更新符合条件的日志：如果 StartTime 和 EndTime 不相等，则更新 EndTime 为当前时间
		for _, orderStatusLog := range statusLogs {
			orderStatusLog.EndTime = currentTime
			if err = DB.Save(&orderStatusLog).Error; err != nil {
				return nil, err
			}
		}

		// 创建并记录取消状态的日志
		statusLog := &models.OrderStatusLog{
			OrderDetailId: int(detail.ID),
			Status:        6,
			StartTime:     currentTime,
			Description:   req.CancelReason,
		}

		if err = DB.Create(&statusLog).Error; err != nil {
			return nil, err
		}
	}

	return &order_common.Empty{}, nil
}

// Reminder implements the OrderUserServiceImpl interface.
// 提醒商家发货
// TODO 需要修改
// 步骤:查询订单是否属于用户->判断status==2,不是就返回错误,是继续(增加健壮性)->通过订单查询商家id,查不到就报错,否则继续->发送信息
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {

	// 查询该用户的订单信息
	var order models.Order

	if err = DB.Where("id = ? AND user_id = ?", req.OrderId, req.UserId).First(&order).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	var orderDetail models.OrderDetail

	if err = DB.Where("order_id = ?", req.OrderId).First(&orderDetail).Error; err != nil {
		return nil, err
	}

	var product models.Product

	if err = DB.Where("id = ?", orderDetail.ProductId).First(&product).Error; err != nil {
		return nil, err
	}

	//提醒商家发货逻辑
	util.SendMessage("order.direct", "message", "orderId:"+strconv.Itoa(int(req.OrderId))+",shopId"+strconv.Itoa(int(product.ShopId)), 1)

	return &order_common.Empty{}, nil
}

// Complete implements the OrderUserServiceImpl interface.
// 确认收货
// TODO 需要修改
// 步骤:查询要完成的订单->查询不到，返回错误，否则继续->判断该订单的status是否为4,否返回错误,是继续->修改订单日志存储状态，如果报错就返回错误->返回正确响应
func (s *OrderUserServiceImpl) Complete(ctx context.Context, req *userOrder.CompleteReq) (resp *order_common.Empty, err error) {

	// 查询订单信息
	var order models.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	if orderResult.Error != nil {
		return nil, orderResult.Error
	}

	// 更新订单的状态为“已完成”
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails)

	if detailResult.Error != nil {
		return nil, detailResult.Error
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	for _, detail := range orderDetails {

		// 查找并更新所有状态为 4 的日志记录
		var statusLogs []models.OrderStatusLog

		if err = DB.Where("order_detail_id = ? AND status = ?", detail.ID, 4).Find(&statusLogs).Error; err != nil {
			log.Println(err)
			return nil, SearchOrderLogError
		}

		// 更新所有状态为4的日志的结束时间
		for _, orderStatusLog := range statusLogs {
			orderStatusLog.EndTime = currentTime
			if err = DB.Save(&orderStatusLog).Error; err != nil {
				log.Println(err)
				return nil, SaveLogError
			}
		}

		// 更新订单详情状态为“已完成”

		if err = DB.Save(&detail).Error; err != nil {
			return nil, err
		}

		// 记录新的状态变更到 OrderStatusLog
		statusLog := &models.OrderStatusLog{
			OrderDetailId: int(detail.ID),
			Status:        5, // 已完成
			StartTime:     currentTime,
			EndTime:       currentTime,
			Description:   "订单已完成",
		}

		if err = DB.Create(&statusLog).Error; err != nil {
			return nil, err
		}
	}

	return &order_common.Empty{}, nil
}
