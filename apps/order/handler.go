package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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

// Submit implements the OrderUserServiceImpl interface.
// 用户提交订单
func (s *OrderUserServiceImpl) Submit(ctx context.Context, req *userOrder.OrderSubmitReq) (resp *userOrder.OrderSubmitResp, err error) {

	orderUUID := uuid.New().String() // 使用 UUID 生成唯一的订单号

	// 创建订单对象
	order := models.Order{
		Number:        orderUUID,
		UserId:        req.UserId,
		AddressBookId: int(req.AddressBookId),
		PayMethod:     int(req.PayMethod),
		Remark:        req.Remark,
		Amount:        float64(req.Amount),
	}
	err = DB.Create(&order).Error
	if err != nil {
		return nil, err
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
			Status:    0, // 初始状态为待付款
		}

		// 将订单详情添加到列表中
		orderDetails = append(orderDetails, orderDetail)
	}

	err = DB.Create(&orderDetails).Error
	if err != nil {
		return nil, err
	}

	for _, detail := range orderDetails {
		orderStatusLog := models.OrderStatusLog{
			OrderDetailId: int(detail.ID),
			Status:        0,
			StartTime:     time.Now(),
			EndTime:       time.Now(),
			Description:   "订单创建，待付款",
		}

		// 保存状态日志
		err = DB.Create(&orderStatusLog).Error
		if err != nil {
			return nil, err
		}
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
	result := DB.Where("user_id = ?", req.UserId).
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	// 构建订单响应数据
	orderList := make([]*order_common.OrderResp, len(orders))
	for i, order := range orders {
		// 根据订单号查询订单详情
		var orderDetails []models.OrderDetail
		DB.Where("order_id = ?", order.ID).Find(&orderDetails)

		// 构建订单响应
		orderResp := &order_common.OrderResp{
			Order: &order_common.Order{
				UserId: order.UserId,
				Number: order.Number,
				Status: int32(order.PayStatus),
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

	// 查询订单信息（Order）基本信息
	var order models.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	// 如果没有找到对应订单，返回错误信息
	if orderResult.Error != nil {
		return nil, orderResult.Error
	}

	// 查询订单详情（List）
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails)

	// 如果查询订单详情失败，返回错误
	if detailResult.Error != nil {
		return nil, detailResult.Error
	}

	// 处理类型
	orderCommon := &order_common.Order{
		Number:        order.Number,
		UserId:        order.UserId,
		PayMethod:     int32(order.PayMethod),
		Status:        int32(order.PayStatus),
		AddressBookId: int32(order.AddressBookId),
		Amount:        float32(order.Amount),
		Remark:        order.Remark,
		Phone:         order.Phone,
		Address:       order.Address,
		Username:      order.UserName,
		Consignee:     order.Consignee,
	}

	// 构建并返回响应数据
	resp = &order_common.OrderResp{
		Order: orderCommon,
	}

	return resp, nil
}

// Cancel implements the OrderUserServiceImpl interface.
// 取消订单
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	// 查询订单信息
	var order models.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	if orderResult.Error != nil {
		return nil, orderResult.Error
	}

	// 查询订单详情（OrderDetails）
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails)

	if detailResult.Error != nil {
		return nil, detailResult.Error
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	// 更新订单详情的状态为“已取消”并记录到 OrderStatusLog
	for _, detail := range orderDetails {
		// 更新订单详情的状态为“已取消”
		detail.Status = 6
		if err := DB.Save(&detail).Error; err != nil {
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
			if err := DB.Save(&orderStatusLog).Error; err != nil {
				return nil, err
			}
		}

		// 创建并记录取消状态的日志
		statusLog := &models.OrderStatusLog{
			OrderDetailId: int(detail.ID),
			Status:        6,
			StartTime:     currentTime,
			EndTime:       currentTime,
			Description:   req.CancelReason,
		}

		if err := DB.Create(&statusLog).Error; err != nil {
			return nil, err
		}
	}

	return &order_common.Empty{}, nil
}

// Reminder implements the OrderUserServiceImpl interface.
// 提醒商家发货
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {

	// 查询该用户的订单信息
	var order models.Order
	orderResult := DB.Where("id = ? AND user_id = ?", req.OrderId, req.UserId).First(&order)

	if orderResult.Error != nil {
		return nil, orderResult.Error
	}

	// TODO 提醒商家发货逻辑

	return &order_common.Empty{}, nil
}

// Complete implements the OrderUserServiceImpl interface.
// 确认收货
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
		if detail.Status == 4 {
			// 查找并更新所有状态为 4 的日志记录
			var statusLogs []models.OrderStatusLog
			statusLogResult := DB.Where("order_detail_id = ? AND status = ?", detail.ID, 4).Find(&statusLogs)
			if statusLogResult.Error != nil {
				return nil, statusLogResult.Error
			}

			// 更新所有状态为4的日志的结束时间
			for _, orderStatusLog := range statusLogs {
				orderStatusLog.EndTime = currentTime
				if err := DB.Save(&orderStatusLog).Error; err != nil {
					return nil, err
				}
			}
		}

		// 更新订单详情状态为“已完成”
		detail.Status = 5
		if err := DB.Save(&detail).Error; err != nil {
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

		if err := DB.Create(&statusLog).Error; err != nil {
			return nil, err
		}
	}

	return &order_common.Empty{}, nil
}
