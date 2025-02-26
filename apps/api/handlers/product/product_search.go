package product

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// Search 搜索商品
func Search(ctx context.Context, c *app.RequestContext) {
	query := c.Query("query")
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "page参数错误",
		})
		return
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "pageSize参数错误",
		})
		return
	}

	products, err := client.SearchProducts(ctx, query, page, pageSize)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"products": products,
	})
}
