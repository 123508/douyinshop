package payment

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Notify(ctx context.Context, c *app.RequestContext) {
	transactionId := c.PostForm("trade_no")
	orderId := c.PostForm("out_trade_no")
	tradeStatus := c.PostForm("trade_status")
	// 更新订单状态等业务逻辑
	if tradeStatus == "TRADE_SUCCESS" {
		client.Notify(ctx, orderId, transactionId)
	}
	// 返回 ACK 确认
	c.String(consts.StatusOK, "success")
}
