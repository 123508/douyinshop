package shop

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {
	shopID, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "shopID参数错误",
		})
		return
	}
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
	req := shop.GetProductListReq{
		ShopId:   uint32(shopID),
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}
	resp, err := client.GetProductList(ctx, &req)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	product := make([]utils.H, 0)
	for _, v := range resp.Products {
		categories := make([]string, 0)
		for _, category := range v.Categories {
			categories = append(categories, category)
		}
		product = append(product, utils.H{
			"name":        v.Name,
			"product_id":  v.Id,
			"description": v.Description,
			"picture":     v.Picture,
			"price":       fmt.Sprintf("%.2f", v.Price),
			"categories":  categories,
			"sales":       v.Sales,
		})
	}
	c.JSON(consts.StatusOK, utils.H{
		"products": product,
	})
}
