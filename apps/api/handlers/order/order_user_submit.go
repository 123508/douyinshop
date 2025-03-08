package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

func Submit(ctx context.Context, c *app.RequestContext) {

	// 获取并解析 user_id 参数
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	type ReqParam struct {
		Address_book_id int
		Pay_method      int
		Remark          string
		Amount          string
	}

	param := &ReqParam{}

	if err := c.Bind(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"err": err,
		})
	}

	// 获取并绑定订单信息
	order := order_common.OrderReq{}
	err := c.Bind(&order)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"error": "订单参数错误"})
		return
	}

	float, _ := strconv.ParseFloat(param.Amount, 32)

	// 调用 UserSubmit 函数提交订单
	result, err := client.UserSubmit(ctx, uint(userId), int32(param.Address_book_id), int32(param.Pay_method), param.Remark, float32(float), &order)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, utils.H{"error": "提交订单失败"})
		return
	}

	// 返回订单提交成功的结果
	c.JSON(consts.StatusOK, utils.H{
		"order_id":     result.OrderId,
		"number":       result.Number,
		"order_amount": result.OrderAmount,
	})

}
