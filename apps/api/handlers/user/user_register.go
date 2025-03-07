package user

import (
	"context"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strings"

	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func Register(ctx context.Context, c *app.RequestContext) {
	req := &user.RegisterReq{}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Error binding JSON",
		})
		return
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		split := strings.Split(err.Error(), ":")
		if split[len(split)-1] == " 密码不匹配,请重新输入" {
			c.JSON(consts.StatusBadRequest, utils.H{
				"error": "密码不匹配,请重新输入",
			})
		} else {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"ok": resp,
	})
}
