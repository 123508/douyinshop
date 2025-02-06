package auth

import (
	"context"
	// "log"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func VerifyToken(ctx context.Context, c *app.RequestContext) {
	req := &auth.VerifyTokenReq{
		Token: c.Query("token"),
	}

	resp, err := client.VerifyToken(context.Background(), req)
	if err != nil {
		// log.Fatal(err)
		c.JSON(500, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, utils.H{
		"res": resp,
	})
}
