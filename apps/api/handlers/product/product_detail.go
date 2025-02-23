package product

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// Detail 获取商品详情
func Detail(ctx context.Context, c *app.RequestContext) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "product_id参数错误",
		})
		return
	}

	result, err := client.GetProductDetail(ctx, productId)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"product": result,
	})
}
