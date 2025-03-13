package user

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func Delete(ctx context.Context, c *app.RequestContext) {
	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	req := &user.DeleteReq{UserId: userId}

	bool, err := client.Delete(ctx, req)

	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"ok": bool,
	})
}
