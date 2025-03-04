package product

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// Detail 获取商品详情
func Detail(ctx context.Context, c *app.RequestContext) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "product_id参数错误",
		})
		return
	}

	result, err := client.GetProductDetail(ctx, productId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"product": result,
	})
}
