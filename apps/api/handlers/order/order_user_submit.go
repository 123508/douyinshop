package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func Submit(ctx context.Context, c *app.RequestContext) {

	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": "user_id参数错误"})
		return
	}

	addressBookId, err := strconv.Atoi(c.Param("address_book_id"))
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": "address_book_id参数错误"})
		return
	}

	payMethod, err := strconv.Atoi(c.Param("pay_method"))
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": "pay_method参数错误"})
		return
	}

	remark := c.Param("remark")
	if remark == "" {
		remark = "无备注"
	}

	amount, err := strconv.ParseFloat(c.Param("amount"), 32)
	if err != nil || amount <= 0 {
		c.JSON(400, map[string]interface{}{"error": "amount参数错误"})
		return
	}

	order := order_common.OrderReq{}
	err = c.Bind(&order)
	if err != nil {
		c.JSON(400, map[string]interface{}{"error": "订单参数错误"})
		return
	}

	result, err := client.UserSubmit(ctx, uint(userId), int32(addressBookId), int32(payMethod), remark, float32(amount), &order)
	if err != nil {
		c.JSON(500, map[string]interface{}{"error": "提交订单失败"})
		return
	}

	c.JSON(200, map[string]interface{}{
		"order_id":     result.OrderId,
		"number":       result.Number,
		"order_amount": result.OrderAmount,
	})

}
