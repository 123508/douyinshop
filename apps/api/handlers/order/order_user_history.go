package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func History(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.ParseUint(c.Query("userId"), 10, 32)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "userId 参数错误",
		})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		c.JSON(400, map[string]interface{}{
			"error": "page 参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize < 1 {
		c.JSON(400, map[string]interface{}{
			"error": "pageSize 参数错误",
		})
		return
	}

	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "status 参数错误",
		})
		return
	}

	historyResp, err := client.UserHistory(ctx, uint32(userId), uint32(page), uint32(pageSize), int32(status))
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"total":     historyResp.Total,
		"list":      historyResp.List,
		"page":      historyResp.Page,
		"page_size": historyResp.PageSize,
	})

}
