package cart

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/checkout"
	"github.com/cloudwego/hertz/pkg/app"
)

func Checkout(ctx context.Context, c *app.RequestContext) {
	var err error
	req := &checkout.CheckoutReq{}
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	_, err = client.Checkout(ctx, req)

	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"ok": true,
	})
}
