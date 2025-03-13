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

func GetInfo(ctx context.Context, c *app.RequestContext) {
	//获取用户id
	value, exists := c.Get("userId")

	userId, ok := value.(uint32)

	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	req := &user.GetUserInfoReq{UserId: userId}

	info, err := client.GetUserInfo(ctx, req)

	if err != nil {
		errorno.DealWithError(err, c)
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
