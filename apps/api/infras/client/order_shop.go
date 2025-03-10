package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/order/businessOrder"
	"github.com/123508/douyinshop/kitex_gen/order/businessOrder/orderbusinessservice"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/config"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var orderShopClient orderbusinessservice.Client

func initOrderShopRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := orderbusinessservice.NewClient(
		config.Conf.BusinessOrderConfig.ServiceName,       // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	orderShopClient = c
}

// GetOrderList 商家查询订单列表
// ShopId 商家ID
// Page 页码
// PageSize 每页数量
// 返回商家订单列表 GetOrderListResp
func GetOrderList(ctx context.Context, ShopId uint32, Page uint32, PageSize uint32) (*businessOrder.GetOrderListResp, error) {

	req := &businessOrder.GetOrderListReq{
		ShopId:   ShopId,
		Page:     Page,
		PageSize: PageSize,
	}

	resp, err := orderShopClient.GetOrderList(ctx, req)

	if err != nil {
		return nil, err
	}

	orderListResp := &businessOrder.GetOrderListResp{
		List: make([]*order_common.OrderResp, len(resp.List)),
	}

	for i, order := range resp.List {
		orderListResp.List[i] = &order_common.OrderResp{
			Order: &order_common.Order{
				UserId:        order.Order.UserId,
				Number:        order.Order.Number,
				Status:        order.Order.Status,
				PayMethod:     order.Order.PayMethod,
				Amount:        order.Order.Amount,
				AddressBookId: order.Order.AddressBookId,
				Remark:        order.Order.Remark,
				Phone:         order.Order.Phone,
				Address:       order.Order.Address,
				Username:      order.Order.Username,
				Consignee:     order.Order.Consignee,
				ShopId:        order.Order.ShopId,
			},
		}
	}

	return orderListResp, nil
}

// ShopDetail 商家查询订单详情
// OrderId 订单ID
// List 订单详情列表
// 返回订单详情 OrderResp
func ShopDetail(ctx context.Context, OrderId uint32, List []order_common.OrderDetail) (*order_common.OrderResp, error) {

	var orderDetails []*order_common.OrderDetail
	for i := range List {
		orderDetails = append(orderDetails, &order_common.OrderDetail{
			Name:      List[i].Name,
			Image:     List[i].Image,
			OrderId:   OrderId,
			ProductId: List[i].ProductId,
			Number:    List[i].Number,
			Amount:    List[i].Amount,
		})
	}

	req := &order_common.OrderReq{
		OrderId: OrderId,
		List:    orderDetails,
	}

	resp, err := orderShopClient.Detail(ctx, req)
	if err != nil {
		return nil, err
	}

	orderResp := &order_common.OrderResp{
		Order: &order_common.Order{
			Number:        resp.Order.Number,
			UserId:        resp.Order.UserId,
			PayMethod:     resp.Order.PayMethod,
			Status:        resp.Order.Status,
			AddressBookId: resp.Order.AddressBookId,
			Amount:        resp.Order.Amount,
			Remark:        resp.Order.Remark,
			Phone:         resp.Order.Phone,
			Address:       resp.Order.Address,
			Username:      resp.Order.Username,
			Consignee:     resp.Order.Consignee,
		},
	}

	return orderResp, nil
}

// ShopConfirm 商家确认订单
// OrderId 订单ID
// Status 订单状态
// 返回空结构体 Empty，表示操作成功
func ShopConfirm(ctx context.Context, OrderId uint32, Status int32) (*order_common.Empty, error) {

	req := &businessOrder.ConfirmReq{
		OrderId: OrderId,
		Status:  Status,
	}

	resp, err := orderShopClient.Confirm(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &order_common.Empty{}, nil
	}

	return resp, nil

}

// ShopDelivery 商家发货
// OrderId 订单ID
// 返回空结构体 Empty，表示操作成功
func ShopDelivery(ctx context.Context, OrderId uint32) (*order_common.Empty, error) {

	req := &businessOrder.DeliveryReq{
		OrderId: OrderId,
	}

	resp, err := orderShopClient.Delivery(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &order_common.Empty{}, nil
	}

	return resp, nil

}

// ShopReceive 商家确认收货
// OrderId 订单ID
// 返回空结构体 Empty，表示操作成功
func ShopReceive(ctx context.Context, OrderId uint32) (*order_common.Empty, error) {

	req := &businessOrder.ReceiveReq{
		OrderId: OrderId,
	}

	resp, err := orderShopClient.Receive(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &order_common.Empty{}, nil
	}

	return resp, nil
}

// ShopRejection 商家拒绝订单
// OrderId 订单ID
// RejectionReason 拒绝原因
// 返回空结构体 Empty，表示操作成功
func ShopRejection(ctx context.Context, OrderId uint32, RejectionReason string) (*order_common.Empty, error) {

	req := &businessOrder.RejectionReq{
		OrderId:         OrderId,
		RejectionReason: RejectionReason,
	}

	resp, err := orderShopClient.Rejection(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &order_common.Empty{}, nil
	}

	return resp, nil

}

// ShopCancel 商家取消订单
// OrderId 订单ID
// CancelReason 取消原因
// 返回空结构体 Empty，表示操作成功
func ShopCancel(ctx context.Context, OrderId uint32, CancelReason string) (*order_common.Empty, error) {

	req := &order_common.CancelReq{
		OrderId:      OrderId,
		CancelReason: CancelReason,
	}

	resp, err := orderShopClient.Cancel(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return &order_common.Empty{}, nil
	}

	return resp, nil

}
