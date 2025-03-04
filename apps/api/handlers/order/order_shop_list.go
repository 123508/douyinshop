package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {

	shopId, err := strconv.Atoi(c.Query("shopId"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "shopId 参数错误",
		})
		return
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "page 参数错误",
		})
		return
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "pageSize 参数错误",
		})
		return
	}

	orderResp, err := client.GetOrderList(ctx, uint32(shopId), uint32(page), uint32(pageSize))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"orders": orderResp,
	})
}
