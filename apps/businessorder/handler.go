package main

import (
	"context"
	businessOrder "github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
)

// OrderBusinessServiceImpl implements the last service interface defined in the IDL.
type OrderBusinessServiceImpl struct{}

// GetOrderList implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) GetOrderList(ctx context.Context, req *businessOrder.GetOrderListReq) (resp *businessOrder.GetOrderListResp, err error) {
	// TODO: Your code here...
	return
}

// Detail implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// Confirm implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Confirm(ctx context.Context, req *businessOrder.ConfirmReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Delivery implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Delivery(ctx context.Context, req *businessOrder.DeliveryReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Receive implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Receive(ctx context.Context, req *businessOrder.ReceiveReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Rejection implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Rejection(ctx context.Context, req *businessOrder.RejectionReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Cancel implements the OrderBusinessServiceImpl interface.
func (s *OrderBusinessServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}
