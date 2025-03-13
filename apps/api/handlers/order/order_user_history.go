package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func History(ctx context.Context, c *app.RequestContext) {

	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	type Param struct {
		Page     uint32
		PageSize uint32
		Status   int32
	}

	param := &Param{}

	if err := c.BindJSON(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "参数错误",
		})
		return
	}

	historyResp, err := client.UserHistory(ctx, userId, param.Page, param.PageSize, param.Status)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"total":     historyResp.Total,
		"list":      historyResp.List,
		"page":      historyResp.Page,
		"page_size": historyResp.PageSize,
	})

}
