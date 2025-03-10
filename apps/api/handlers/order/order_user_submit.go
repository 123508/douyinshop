package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/errorno"
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

	type Detail struct {
		ProductId uint32 `json:"product_id"`
		Name      string
		Image     string
		Number    uint32
		Amount    string
	}

	type ReqParam struct {
		AddressBookId int `json:"address_book_id"`
		PayMethod     int `json:"pay_method"`
		Remark        string
		Amount        string
		List          []Detail
	}

	param := &ReqParam{}

	if err := c.Bind(param); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"err": err,
		})
	}

	// 获取并绑定订单信息
	order := order_common.OrderReq{}

	for _, k := range param.List {
		float, err := strconv.ParseFloat(k.Amount, 32)
		if err != nil {
			c.JSON(400, utils.H{
				"err": "参数错误",
			})
		}

		order.List = append(order.List, &order_common.OrderDetail{
			Name:      k.Name,
			Image:     k.Image,
			OrderId:   0,
			ProductId: k.ProductId,
			Number:    k.Number,
			Amount:    float32(float),
		})
	}

	float, err := strconv.ParseFloat(param.Amount, 32)

	if err != nil {
		c.JSON(400, utils.H{
			"err": "参数错误",
		})
	}

	// 调用 UserSubmit 函数提交订单
	result, err := client.UserSubmit(ctx, userId, int32(param.AddressBookId), int32(param.PayMethod), param.Remark, float32(float), &order)
	if err != nil {
		basicErr := errorno.ParseBasicMessageError(err)

		if basicErr.Raw != nil {
			c.JSON(consts.StatusInternalServerError, utils.H{
				"err": err,
			})
		} else {
			c.JSON(basicErr.Code, utils.H{
				"error": basicErr.Message,
			})
		}
	}

	// 返回订单提交成功的结果
	c.JSON(consts.StatusOK, utils.H{
		"order_id":     result.OrderId,
		"number":       result.Number,
		"order_amount": result.OrderAmount,
	})

}
