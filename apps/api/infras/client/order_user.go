package client

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder"
	"github.com/123508/douyinshop/kitex_gen/order/userOrder/orderuserservice"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/models"

	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var orderUserClient orderuserservice.Client

func initOrderUserRpc() {
	r, err := etcd.NewEtcdResolverWithAuth(config.Conf.EtcdConfig.Endpoints, config.Conf.EtcdConfig.Username, config.Conf.EtcdConfig.Password)
	if err != nil {
		panic(err)
	}

	c, err := orderuserservice.NewClient(
		config.Conf.OrderConfig.ServiceName,               // service name
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	orderUserClient = c
}

// UserSubmit 用户提交订单
// userId 用户id
// addressBookId 地址id
// payMethod 支付方式
// remark 用户备注
// amount 产品数量
// order 订单信息
// 返回订单提交resp
func UserSubmit(ctx context.Context, userId uint, addressBookId int32, payMethod int32, remark string, amount float32, order *order_common.OrderReq) (*userOrder.OrderSubmitResp, error) {

	req := &userOrder.OrderSubmitReq{
		UserId:        uint32(userId),
		AddressBookId: addressBookId,
		PayMethod:     payMethod,
		Remark:        remark,
		Amount:        amount,
		Order:         order,
	}

	resp, err := orderUserClient.Submit(ctx, req)

	if err != nil {
		return nil, err
	}

	result := &userOrder.OrderSubmitResp{
		OrderId:     resp.OrderId,
		Number:      resp.Number,
		OrderAmount: resp.OrderAmount,
	}

	return result, nil
}

// UserHistory 用户查询历史订单
// userId 用户ID
// page 页码
// pageSize 每页数量
// status 订单状态
// 返回用户历史订单列表 HistoryResp
func UserHistory(ctx context.Context, userId uint32, page uint32, pageSize uint32, status int32) (*userOrder.HistoryResp, error) {

	req := &userOrder.HistoryReq{
		UserId:   userId,
		Page:     page,
		PageSize: pageSize,
		Status:   status,
	}

	resp, err := orderUserClient.History(ctx, req)

	if err != nil {
		return nil, err
	}

	historyResp := &userOrder.HistoryResp{
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
		List:     make([]*order_common.OrderResp, len(resp.List)),
	}
	for i, order := range resp.List {
		historyResp.List[i] = &order_common.OrderResp{
			Order: &order_common.Order{
				UserId: order.Order.UserId,
				Number: order.Order.Number,
				Status: order.Order.Status,
			},
		}
	}

	return historyResp, nil

}

// UserDetail 用户查询订单详情
// orderId 订单ID
// List 订单详细信息列表
// 返回订单详情 OrderResp
func UserDetail(ctx context.Context, orderId uint32, List []models.OrderDetail) (*order_common.OrderResp, error) {

	var orderDetails []*order_common.OrderDetail
	for _, detail := range List {
		orderDetails = append(orderDetails, &order_common.OrderDetail{
			Name:      detail.Name,
			Image:     detail.Image,
			OrderId:   orderId,
			ProductId: uint32(detail.ProductId),
			Number:    uint32(detail.Number),
			Amount:    float32(detail.Amount),
		})
	}

	req := &order_common.OrderReq{
		OrderId: orderId,
		List:    orderDetails,
	}

	// 调用订单服务的详情查询接口
	resp, err := orderUserClient.Detail(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建返回的订单详情响应
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

// UserCancel 用户取消订单
// orderId 订单ID
// CancelReason 取消原因
// 返回空结构体 Empty 表示取消结果
func UserCancel(ctx context.Context, orderId uint32, CancelReason string) (order_common.Empty, error) {

	req := &order_common.CancelReq{
		OrderId:      orderId,
		CancelReason: CancelReason,
	}

	_, err := orderUserClient.Cancel(ctx, req)
	if err != nil {
		return order_common.Empty{}, err
	}

	return order_common.Empty{}, nil
}

// UserReminder 用户提醒商家发货
// userid 用户ID
// orderId 订单ID
// 返回空结构体 Empty 表示提醒成功
func UserReminder(ctx context.Context, userid uint32, orderId uint32) (order_common.Empty, error) {

	req := &userOrder.ReminderReq{
		UserId:  userid,
		OrderId: orderId,
	}

	_, err := orderUserClient.Reminder(ctx, req)
	if err != nil {
		return order_common.Empty{}, err
	}

	return order_common.Empty{}, nil
}

// UserComplete 用户确认收货
// orderId 订单ID
// 返回空结构体 Empty 表示确认收货成功
func UserComplete(ctx context.Context, orderId uint32) (order_common.Empty, error) {

	req := &userOrder.CompleteReq{
		OrderId: orderId,
	}

	_, err := orderUserClient.Complete(ctx, req)
	if err != nil {
		return order_common.Empty{}, err
	}

	return order_common.Empty{}, nil
}
