package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Cancel(ctx context.Context, c *app.RequestContext) {

	type Param struct {
		CancelReason string `json:"cancel_reason"`
		OrderId      uint32 `json:"order_id"`
	}

	param := &Param{}

	if err := c.BindJSON(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	resp, err := client.UserCancel(ctx, param.OrderId, param.CancelReason)
	if err != nil {
		basicErr := errorno.ParseBasicMessageError(err)

		if basicErr.Raw != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"err": err,
			})
		} else {
			c.JSON(basicErr.Code, utils.H{
				"error": basicErr.Message,
			})
		}
	}

	c.JSON(consts.StatusOK, utils.H{
		"message": "订单已取消",
		"result":  &resp,
	})

}
