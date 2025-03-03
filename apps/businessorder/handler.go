package main

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/db"
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

// GetOrderList implements the OrderBusinessServiceImpl interface.
// 获取订单
func (s *OrderBusinessServiceImpl) GetOrderList(ctx context.Context, req *businessOrder.GetOrderListReq) (resp *businessOrder.GetOrderListResp, err error) {

	// 查询产品ID列表（根据 shopId）
	var products []models.Product
	productResult := DB.Where("shop_id = ?", req.ShopId).Find(&products)
	if productResult.Error != nil {
		return nil, productResult.Error
	}

	// 获取所有产品ID
	var productIds []uint32
	for _, product := range products {
		productIds = append(productIds, uint32(product.ID))
	}

	// 查询订单详情（OrderDetail），通过产品ID
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("product_id IN ?", productIds).Find(&orderDetails)
	if detailResult.Error != nil {
		return nil, detailResult.Error
	}

	// 查询与订单详情相关的所有订单（Order）
	var orderIds []uint32
	for _, detail := range orderDetails {
		orderIds = append(orderIds, uint32(detail.OrderId))
	}

	var orders []models.Order
	orderResult := DB.Where("id IN ?", orderIds).Find(&orders)
	if orderResult.Error != nil {
		return nil, orderResult.Error
	}

	// 封装订单响应列表
	var orderRespList []*order_common.OrderResp
	for _, order := range orders {

		// 将 models.Order 转换为 order_common.Order
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

// Confirm implements the OrderBusinessServiceImpl interface.
// 确认订单
func (s *OrderBusinessServiceImpl) Confirm(ctx context.Context, req *businessOrder.ConfirmReq) (resp *order_common.Empty, err error) {

	if err := s.updateOrderStatus(req.OrderId, 2); err != nil {
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
	if err := s.updateOrderStatus(req.OrderId, 4); err != nil {
		return nil, err
	}

	// 返回空响应
	return &order_common.Empty{}, nil
}

// Rejection implements the OrderBusinessServiceImpl interface.
// 商家拒绝订单
func (s *OrderBusinessServiceImpl) Rejection(ctx context.Context, req *businessOrder.RejectionReq) (resp *order_common.Empty, err error) {

	// 查询订单详情（OrderDetails），根据 OrderId 获取订单详情
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("order_id = ?", req.OrderId).Find(&orderDetails)

	if detailResult.Error != nil {
		return nil, detailResult.Error
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	// 更新当前订单详情的状态为“拒单”
	for _, detail := range orderDetails {
		// 只更新状态为拒单（状态9）
		if detail.Status != 9 {
			// 更新订单详情的状态为拒单（状态 9）
			detail.Status = 9
			if err := DB.Save(&detail).Error; err != nil {
				return nil, err
			}

			// 查找并更新与 OrderDetailId 匹配的日志的 EndTime
			var statusLogs []models.OrderStatusLog
			statusLogResult := DB.Where("order_detail_id = ? AND status != ?", detail.ID, 9).Find(&statusLogs)

			if statusLogResult.Error != nil {
				return nil, statusLogResult.Error
			}

			// 更新状态日志的 EndTime，判断 StartTime 和 EndTime 不相等时
			for _, orderStatusLog := range statusLogs {
				if orderStatusLog.StartTime != orderStatusLog.EndTime { // 只有 StartTime 与 EndTime 不相等时才更新
					orderStatusLog.EndTime = currentTime
					if err := DB.Save(&orderStatusLog).Error; err != nil {
						return nil, err
					}
				}
			}

			// 创建并记录拒单状态的日志
			statusLog := &models.OrderStatusLog{
				OrderDetailId: int(detail.ID),
				Status:        9,
				StartTime:     currentTime,
				EndTime:       currentTime,
				Description:   req.RejectionReason,
			}

			if err := DB.Create(&statusLog).Error; err != nil {
				return nil, err
			}
		}
	}

	return &order_common.Empty{}, nil
}

// Cancel implements the OrderBusinessServiceImpl interface.
// 商家取消订单
func (s *OrderBusinessServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {

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

// 核心方法，用于更新订单的状态
func (s *OrderBusinessServiceImpl) updateOrderStatus(orderId uint32, status int) error {

	// 查询订单信息（Order）
	var order order_common.Order
	orderResult := DB.Where("id = ?", orderId).First(&order)

	if orderResult.Error != nil {
		return orderResult.Error
	}

	// 查询订单详情（OrderDetails）
	var orderDetails []models.OrderDetail
	detailResult := DB.Where("order_id = ?", orderId).Find(&orderDetails)

	if detailResult.Error != nil {
		return detailResult.Error
	}

	// 当前时间，用于记录日志
	currentTime := time.Now()

	// 更新所有订单详情的状态
	for _, detail := range orderDetails {
		// 查找并更新所有状态为 status-1 的日志记录的结束时间
		var statusLogs []models.OrderStatusLog
		statusLogResult := DB.Where("order_detail_id = ? AND status = ?", detail.ID, status-1).Find(&statusLogs)
		if statusLogResult.Error != nil {
			return statusLogResult.Error
		}

		// 更新所有状态为 status-1 的日志的结束时间
		for _, orderStatusLog := range statusLogs {
			// 只有 StartTime 和 EndTime 不相等时才更新
			if orderStatusLog.StartTime != orderStatusLog.EndTime {
				orderStatusLog.EndTime = currentTime
				if err := DB.Save(&orderStatusLog).Error; err != nil {
					return err
				}
			}
		}

		// 更新订单详情状态为指定的 status
		detail.Status = status
		if err := DB.Save(&detail).Error; err != nil {
			return err
		}

		// 创建并记录新的状态日志
		statusLog := &models.OrderStatusLog{
			OrderDetailId: int(detail.ID),
			Status:        uint(status),
			StartTime:     currentTime,
			EndTime:       currentTime,
			Description:   "订单状态已更新",
		}

		if err := DB.Create(&statusLog).Error; err != nil {
			return err
		}
	}

	return nil
}
