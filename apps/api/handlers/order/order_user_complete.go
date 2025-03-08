package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Complete(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}
	resp, err := client.UserComplete(ctx, uint32(orderId))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "收货确认成功",
		"result":  &resp,
	})

}
