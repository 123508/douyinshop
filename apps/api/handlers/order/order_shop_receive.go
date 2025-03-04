package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Receive(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}

	resp, err := client.ShopReceive(ctx, uint32(orderId))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "订单确认收货失败",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "订单确认收货成功",
		"data":    resp,
	})

}
