package cart

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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

	req := &cart.GetCartReq{UserId: userId}
	resp, err := client.GetCart(ctx, req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	cartItems := make([]map[string]interface{}, 0)
	for _, item := range resp.Cart.Items {
		cartItems = append(cartItems, map[string]interface{}{
			"product_id":  item.ProductId,
			"product_num": item.Quantity,
		})
	}

	c.JSON(consts.StatusOK, utils.H{
		"cart": cartItems,
	})
}
