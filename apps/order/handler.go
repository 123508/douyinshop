package main

import (
	"context"
	"fmt"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
	userOrder "github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
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

	log.Printf("Received History request for user ID: %d, page: %d, page size: %d", req.UserId, req.Page, req.PageSize)

	// 查询用户订单数据
	var orders []models.Order
	result := DB.Where("user_id = ?", req.UserId).
		Offset(int(offset)).
		Limit(int(req.PageSize)).
		Find(&orders)

	if result.Error != nil {
		log.Printf("Error fetching order history for user ID %d: %v", req.UserId, result.Error)
		return nil, fmt.Errorf("failed to fetch order history: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Printf("No orders found for user ID %d on page %d with page size %d", req.UserId, req.Page, req.PageSize)
		return &userOrder.HistoryResp{List: []*order_common.OrderResp{}}, nil
	}

	// 构建订单响应数据
	orderList := make([]*order_common.OrderResp, len(orders))
	for i, order := range orders {
		// 根据订单号查询订单详情
		var orderDetails []models.OrderDetail
		DB.Where("order_number = ?", order.Number).Find(&orderDetails)

		// 构建订单响应
		orderResp := &order_common.OrderResp{
			Order: &order_common.Order{
				UserId: order.UserId,
				Number: order.Number, // 添加订单号字段

			},
		}

		// 处理订单详情
		for _, detail := range orderDetails {
			orderResp.List = append(orderResp.List, &order_common.OrderDetail{
				OrderId: uint32(detail.OrderId),
			})
		}
		orderList[i] = orderResp
	}

	// 返回响应
	return &userOrder.HistoryResp{
		List:     orderList,
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    uint32(len(orders)), // 假设返回的总数是查询到的订单数
	}, nil
}

// Detail implements the OrderUserServiceImpl interface.
// 查询订单的详细信息
func (s *OrderUserServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {

	// 输入参数验证，确保传递的 orderId 合法
	if req.OrderId <= 0 {
		return nil, fmt.Errorf("无效的订单ID")
	}

	// 查询订单信息（Order）基本信息
	var order order_common.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	// 如果没有找到对应订单，返回错误信息
	if orderResult.Error != nil {
		return nil, fmt.Errorf("查询订单信息失败：%w", orderResult.Error)
	}

	// 查询订单详情（List）
	var orderDetails []order_common.OrderDetail
	detailResult := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails)

	// 如果查询订单详情失败，返回错误
	if detailResult.Error != nil {
		return nil, fmt.Errorf("查询订单详情失败：%w", detailResult.Error)
	}

	// 将 orderDetails 转换为指针切片 []*order_common.OrderDetail
	var orderDetailPtrs []*order_common.OrderDetail
	for i := range orderDetails {
		orderDetailPtrs = append(orderDetailPtrs, &orderDetails[i])
	}

	// 构建并返回响应数据
	resp = &order_common.OrderResp{
		Order: &order,          // 返回订单的基本信息
		List:  orderDetailPtrs, // 返回订单详情列表（指针类型）
	}

	return resp, nil
}

// Cancel implements the OrderUserServiceImpl interface.
// 取消订单
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

	// 输入参数验证，确保传递的 orderId 合法
	if req.OrderId <= 0 {
		return nil, fmt.Errorf("无效的订单ID")
	}

	// 查询订单信息
	var order order_common.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	if orderResult.Error != nil {
		return nil, fmt.Errorf("查询订单信息失败：%w", orderResult.Error)
	}

	// 更新订单的状态为“已取消”
	order.Status = 7
	order.CancelReason = req.CancelReason
	if err := DB.Save(&order).Error; err != nil {
		return nil, fmt.Errorf("更新订单状态失败：%w", err)
	}

	return &order_common.Empty{}, nil
}

// Reminder implements the OrderUserServiceImpl interface.
// 提醒商家发货
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {

	// 输入参数验证，确保传递的 orderId 和 userId 合法
	if req.OrderId <= 0 || req.UserId <= 0 {
		return nil, fmt.Errorf("无效的订单ID或用户ID")
	}

	// 查询该用户的订单信息
	var order order_common.Order
	orderResult := DB.Where("id = ? AND user_id = ?", req.OrderId, req.UserId).First(&order)

	// 如果没有找到对应订单，返回错误信息
	if orderResult.Error != nil {
		return nil, fmt.Errorf("未找到该用户的订单，订单ID: %d, 用户ID: %d", req.OrderId, req.UserId)
	}

	// 检查订单状态，假设只有未发货的订单可以进行提醒
	if order.Status != 2 {
		return nil, fmt.Errorf("该订单无法提醒发货，订单状态: %s", order.Status)
	}

	// TODO 提醒商家发货逻辑

	// 返回空响应，表示操作成功
	return &order_common.Empty{}, nil
}

// Complete implements the OrderUserServiceImpl interface.
// 确认收货
func (s *OrderUserServiceImpl) Complete(ctx context.Context, req *userOrder.CompleteReq) (resp *order_common.Empty, err error) {

	// 输入参数验证，确保传递的 orderId 合法
	if req.OrderId <= 0 {
		return nil, fmt.Errorf("无效的订单ID")
	}

	// 查询订单信息
	var order order_common.Order
	orderResult := DB.Where("id = ?", req.OrderId).First(&order)

	if orderResult.Error != nil {
		return nil, fmt.Errorf("查询订单信息失败：%w", orderResult.Error)
	}

	// 更新订单的状态为“已完成”
	order.Status = 6
	if err := DB.Save(&order).Error; err != nil {
		return nil, fmt.Errorf("更新订单状态失败：%w", err)
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}
