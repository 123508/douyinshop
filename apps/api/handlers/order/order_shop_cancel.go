package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func CancelShop(ctx context.Context, c *app.RequestContext) {

	// 获取并解析 user_id 参数
	value, exists := c.Get("userId")
	_, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	type Param struct {
		OrderId      uint32 `json:"order_id"`
		CancelReason string `json:"cancel_reason"`
	}

	param := &Param{}

	if err := c.BindJSON(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	resp, err := client.ShopCancel(ctx, param.OrderId, param.CancelReason)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"message": "订单取消成功",
		"data":    resp,
	})

}
