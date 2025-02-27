package shop

import (
	"context"
	"fmt"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/shop"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func List(ctx context.Context, c *app.RequestContext) {
	shopID, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "shopID参数错误",
		})
		return
	}
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
	req := shop.GetProductListReq{
		ShopId:   uint32(shopID),
		Page:     uint32(page),
		PageSize: uint32(pageSize),
	}
	resp, err := client.GetProductList(ctx, &req)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	product := make([]map[string]interface{}, 0)
	for _, v := range resp.Products {
		categories := make([]string, 0)
		for _, category := range v.Categories {
			categories = append(categories, category)
		}
		product = append(product, map[string]interface{}{
			"name":        v.Name,
			"description": v.Description,
			"picture":     v.Picture,
			"price":       fmt.Sprintf("%.2f", v.Price),
			"categories":  categories,
			"sales":       v.Sales,
		})
	}
	c.JSON(200, map[string]interface{}{
		"products": product,
	})
}
