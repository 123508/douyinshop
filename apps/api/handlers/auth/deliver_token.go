package auth

import (
	"context"
	// "log"
	"strconv"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func DeliverToken(ctx context.Context, c *app.RequestContext) {
	user_id, err := strconv.Atoi(c.Query("userId"))
	req := &auth.DeliverTokenReq{
		UserId: int32(user_id),
	}

	if err != nil {
		c.JSON(400, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	resp, err := client.DeliverToken(context.Background(), req)
	if err != nil {
		// log.Fatal(err)
		c.JSON(500, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, utils.H{
		"token": resp,
	})
}
