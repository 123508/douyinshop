package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"

	"github.com/cloudwego/hertz/pkg/app"
)

func Delete(ctx context.Context, c *app.RequestContext) {
	type request struct {
		AddressID int `json:"addr_id"`
	}
	req := &request{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "参数错误",
		})
		return
	}
	userID, _ := c.Get("userId")
	resp, err := client.DeleteAddress(ctx, req.AddressID, userID.(uint32))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"ok": resp,
	})
}
