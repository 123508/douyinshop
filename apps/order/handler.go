package main

import (
	"context"
	order_common "github.com/123508/douyinshop/kitex_gen/order/order_common"
	userOrder "github.com/123508/douyinshop/kitex_gen/order/userOrder"
)

// OrderUserServiceImpl implements the last service interface defined in the IDL.
type OrderUserServiceImpl struct{}

// Submit implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) Submit(ctx context.Context, req *userOrder.OrderSubmitReq) (resp *userOrder.OrderSubmitResp, err error) {
	// TODO: Your code here...
	return
}

// History implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) History(ctx context.Context, req *userOrder.HistoryReq) (resp *userOrder.HistoryResp, err error) {
	// TODO: Your code here...
	return
}

// Detail implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) Detail(ctx context.Context, req *order_common.OrderReq) (resp *order_common.OrderResp, err error) {
	// TODO: Your code here...
	return
}

// Cancel implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) Cancel(ctx context.Context, req *order_common.CancelReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Reminder implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) Reminder(ctx context.Context, req *userOrder.ReminderReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}

// Complete implements the OrderUserServiceImpl interface.
func (s *OrderUserServiceImpl) Complete(ctx context.Context, req *userOrder.CompleteReq) (resp *order_common.Empty, err error) {
	// TODO: Your code here...
	return
}
