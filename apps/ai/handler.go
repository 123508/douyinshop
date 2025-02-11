package main

import (
	"context"
	ai "github.com/123508/douyinshop/kitex_gen/ai"
)

// AiServiceImpl implements the last service interface defined in the IDL.
type AiServiceImpl struct{}

// OrderQuery implements the AiServiceImpl interface.
func (s *AiServiceImpl) OrderQuery(ctx context.Context, req *ai.OrderQueryReq) (resp *ai.OrderQueryResp, err error) {
	// TODO: Your code here...
	return
}

// AutoPlaceOrder implements the AiServiceImpl interface.
func (s *AiServiceImpl) AutoPlaceOrder(ctx context.Context, req *ai.AutoPlaceOrderReq) (resp *ai.AutoPlaceOrderResp, err error) {
	// TODO: Your code here...
	return
}
