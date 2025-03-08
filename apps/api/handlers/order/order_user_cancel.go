package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Cancel(ctx context.Context, c *app.RequestContext) {

	type Param struct {
		Cancel_reason string
		Order_id      uint32
	}

	param := &Param{}

	if err := c.BindJSON(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	resp, err := client.UserCancel(ctx, param.Order_id, param.Cancel_reason)
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
