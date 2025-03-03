package user

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func GetInfo(ctx context.Context, c *app.RequestContext) {

	value, exists := c.Get("userId")

	fmt.Println(value)

	if !exists {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must exist",
		})
		return
	}

	userId, ok := value.(uint32)

	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
	}

	req := &user.GetUserInfoReq{UserId: userId}

	info, err := client.GetUserInfo(ctx, req)

	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"id":       userId,
		"gender":   info.GetGender(),
		"nickname": info.GetNickname(),
		"phone":    info.GetPhone(),
		"email":    info.GetEmail(),
	})

}
