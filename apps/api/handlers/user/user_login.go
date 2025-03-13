package user

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func Login(ctx context.Context, c *app.RequestContext) {

	req := &user.LoginReq{}

	if err := c.BindJSON(req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err,
		})
		return
	}

	resp, err := client.Login(ctx, req)

	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	token, err := client.DeliverToken(ctx, &auth.DeliverTokenReq{UserId: resp})

	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": err,
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"token": token,
	})

}
