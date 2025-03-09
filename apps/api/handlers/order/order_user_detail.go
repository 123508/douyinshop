package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Detail(ctx context.Context, c *app.RequestContext) {

	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId参数错误",
		})
		return
	}

	orderResp, err := client.UserDetail(ctx, uint32(orderId))
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
		"order": orderResp,
	})

}
