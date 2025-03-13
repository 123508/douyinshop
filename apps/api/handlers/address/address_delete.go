package address

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func Delete(ctx context.Context, c *app.RequestContext) {
	type request struct {
		AddressID int `json:"addr_id"`
	}
	req := &request{}
	err := c.Bind(req)
	if err != nil {
		c.JSON(400, utils.H{
			"error": "参数错误",
		})
		return
	}
	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	resp, err := client.DeleteAddress(ctx, req.AddressID, userId)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"ok": resp,
	})
}
