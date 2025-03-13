package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {
	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	addressList, err := client.GetAddressList(ctx, userId)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"address": addressList,
	})
}
