package main

import (
	"context"
	"fmt"
	businessOrder "github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/db"
	"gorm.io/gorm"
	"log"
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

// GetOrderList implements the OrderBusinessServiceImpl interface.
// 获取订单
func (s *OrderBusinessServiceImpl) GetOrderList(ctx context.Context, req *businessOrder.GetOrderListReq) (resp *businessOrder.GetOrderListResp, err error) {

	return
}

// Detail implements the OrderBusinessServiceImpl interface.
// 订单详情
func (s *OrderBusinessServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {

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

// Confirm implements the OrderBusinessServiceImpl interface.
// 确认订单
func (s *OrderBusinessServiceImpl) Confirm(ctx context.Context, req *businessOrder.ConfirmReq) (resp *order_common.Empty, err error) {

	if err := s.updateOrderStatus(req.OrderId, 3); err != nil {
		return nil, err
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Delivery implements the OrderBusinessServiceImpl interface.
// 商家发货
func (s *OrderBusinessServiceImpl) Delivery(ctx context.Context, req *businessOrder.DeliveryReq) (resp *order_common.Empty, err error) {

	if err := s.updateOrderStatus(req.OrderId, 3); err != nil {
		return nil, err
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Receive implements the OrderBusinessServiceImpl interface.
// 转变订单状态为待收货
func (s *OrderBusinessServiceImpl) Receive(ctx context.Context, req *businessOrder.ReceiveReq) (resp *order_common.Empty, err error) {
	if err := s.updateOrderStatus(req.OrderId, 3); err != nil {
		return nil, err
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Rejection implements the OrderBusinessServiceImpl interface.
// 商家拒单
func (s *OrderBusinessServiceImpl) Rejection(ctx context.Context, req *businessOrder.ReceiveReq) (resp *order_common.Empty, err error) {

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
	// TODO rejectionReason
	order.RejectionReason = ""
	if err := DB.Save(&order).Error; err != nil {
		return nil, fmt.Errorf("更新订单状态失败：%w", err)
	}

	return &order_common.Empty{}, nil
}

// Cancel implements the OrderBusinessServiceImpl interface.
// 商家取消订单
func (s *OrderBusinessServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

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

// 核心方法，用于更新订单的状态
func (s *OrderBusinessServiceImpl) updateOrderStatus(orderId uint32, status int) error {
	if orderId <= 0 {
		return fmt.Errorf("无效的订单ID")
	}

	// 查询订单信息（Order）
	var order order_common.Order
	orderResult := DB.Where("id = ?", orderId).First(&order)

	if orderResult.Error != nil {
		return fmt.Errorf("查询订单信息失败：%w", orderResult.Error)
	}

	// 更新订单的状态
	order.Status = int32(status)
	if err := DB.Save(&order).Error; err != nil {
		return fmt.Errorf("更新订单状态失败：%w", err)
	}

	return nil
}
