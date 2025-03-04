package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Confirm(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}

	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "status 参数错误",
		})
		return
	}

	resp, err := client.ShopConfirm(ctx, uint32(orderId), int32(status))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "订单确认失败",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "订单确认成功",
		"data":    resp,
	})

}
