package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Cancel(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}

	CancelReason := c.Query("cancelReason")
	if CancelReason == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "cancelReason 参数不能为空",
		})
		return
	}

	resp, err := client.UserCancel(ctx, uint32(orderId), CancelReason)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "订单已取消",
		"result":  &resp,
	})

}
