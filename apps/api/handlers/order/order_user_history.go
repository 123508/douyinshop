package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func History(ctx context.Context, c *app.RequestContext) {

	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "page 参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "pageSize 参数错误",
		})
		return
	}

	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "status 参数错误",
		})
		return
	}

	historyResp, err := client.UserHistory(ctx, userId, uint32(page), uint32(pageSize), int32(status))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"total":     historyResp.Total,
		"list":      historyResp.List,
		"page":      historyResp.Page,
		"page_size": historyResp.PageSize,
	})

}
