package user

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app"
)

func UpdateInfo(ctx context.Context, c *app.RequestContext) {
	//获取用户id
	value, exists := c.Get("userId")

	userId, ok := value.(uint32)

	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	req := &user.UpdateReq{}
	//从方法体中填充用户信息
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "Error binding JSON",
		})
		return
	}

	req.UserId = userId
	//发送更新请求
	bool, err := client.Update(ctx, req)

	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"ok": bool,
	})

}
