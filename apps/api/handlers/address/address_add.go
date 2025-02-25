package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/app"
)

func Add(ctx context.Context, c *app.RequestContext) {
	userID, _ := c.Get("userId")
	address := &client.AddressItem{}
	err := c.Bind(address)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "参数错误",
		})
		return
	}
	resp, err := client.AddAddress(ctx, address, userID.(uint32))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"addr_id": resp,
	})
}
