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

	type Req struct {
		ProductId uint32 `json:"product_id"`
		Quantity  uint32 `json:"number"`
	}
	var req Req
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

		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"ok": true,
	})
}
