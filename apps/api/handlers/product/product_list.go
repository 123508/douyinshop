package product

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// List 获取商品列表
func List(ctx context.Context, c *app.RequestContext) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "page参数错误",
		})
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "pageSize参数错误",
		})
		return
	}
	category := c.Query("category")

	products, err := client.ListProducts(ctx, page, pageSize, category)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"products": products,
	})
}
