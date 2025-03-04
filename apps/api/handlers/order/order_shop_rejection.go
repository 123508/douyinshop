package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Rejection(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}

	rejectionReason := c.Query("rejectionReason")
	if rejectionReason == "" {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "拒绝原因不能为空",
		})
		return
	}

	resp, err := client.ShopRejection(ctx, uint32(orderId), rejectionReason)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "订单拒绝失败",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "订单拒绝成功",
		"data":    resp,
	})

}
