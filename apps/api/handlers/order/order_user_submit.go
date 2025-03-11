package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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
		AddressBookId int `json:"address_book_id"`
		PayMethod     int `json:"pay_method"`
		Remark        string
	}

	param := &ReqParam{}

	if err := c.Bind(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"err": err,
		})
	}

	// 调用 UserSubmit 函数提交订单
	result, err := client.UserSubmit(ctx, userId, int32(param.AddressBookId), int32(param.PayMethod), param.Remark)

	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	// 返回订单提交成功的结果
	c.JSON(consts.StatusOK, utils.H{
		"order_id":     result.OrderId,
		"number":       result.Number,
		"order_amount": result.OrderAmount,
	})

}
