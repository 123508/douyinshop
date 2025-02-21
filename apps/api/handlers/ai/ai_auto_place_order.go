package ai

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	ai "github.com/123508/douyinshop/kitex_gen/ai"
	"github.com/123508/douyinshop/kitex_gen/ai/aiservice"
)

// AutoPlaceOrder handles the automatic order placement request
func AutoPlaceOrder(ctx context.Context, c *app.RequestContext) {
	var req ai.AutoPlaceOrderReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}
	
	// 直接创建客户端并调用
	client := aiservice.MustNewClient("ai")
	resp, err := client.AutoPlaceOrder(ctx, &req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code":    500,
			"message": "服务调用失败: " + err.Error(),
		})
		return
	}
	
	c.JSON(200, resp)
}
