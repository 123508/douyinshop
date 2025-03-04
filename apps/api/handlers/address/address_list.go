package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	addressList, err := client.GetAddressList(ctx, userId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"address": addressList,
	})
}
