package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Rejection(ctx context.Context, c *app.RequestContext) {

	type Param struct {
		Order_id         uint32
		Rejection_reason string
	}

	param := &Param{}

	if err := c.BindJSON(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	resp, err := client.ShopRejection(ctx, param.Order_id, param.Rejection_reason)
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
		"message": "订单拒绝成功",
		"data":    resp,
	})

}
