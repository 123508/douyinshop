package cart

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/cart"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Reduce(ctx context.Context, c *app.RequestContext) {
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	var req client.CartReq
	err := c.Bind(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	deleteItemReq := &cart.DeleteItemReq{
		ProductId: req.ProductId,
		Num:       req.Quantity,
		UserId:    userId,
	}

	_, err = client.DeleteItem(ctx, deleteItemReq)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"ok": true,
	})
}
