package user

import (
	"context"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Login(ctx context.Context, c *app.RequestContext) {
	req := &user.LoginReq{
		Email:    c.Query("email"),
		Password: c.Query("password"),
	}

	resp, err := client.Login(context.Background(), req)
	if err != nil {
		// log.Fatal(err)
		c.JSON(500, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, utils.H{
		"userId": resp,
	})
}
