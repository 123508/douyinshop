package user

import (
	"context"
	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Login(ctx context.Context, c *app.RequestContext) {

	req := &user.LoginReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Error binding JSON",
		})
		return
	}

	resp, err := client.Login(ctx, req)

	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	token, err := client.DeliverToken(ctx, &auth.DeliverTokenReq{UserId: resp})

	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": err,
		})
	}

	c.JSON(consts.StatusOK, utils.H{
		"token": token,
	})

}
