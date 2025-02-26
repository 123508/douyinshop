package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {
	userId, _ := c.Get("userId")
	addressList, err := client.GetAddressList(ctx, userId.(uint32))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"address": addressList,
	})
}
