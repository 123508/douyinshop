package payment

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

type Request struct {
	Amount                    string `json:"amount"`
	CreditCardCvv             int    `json:"credit_card_cvv"`
	CreditCardExpirationMonth int    `json:"credit_card_expiration_month"`
	CreditCardExpirationYear  int    `json:"credit_card_expiration_year"`
	CreditCardNumber          string `json:"credit_card_number"`
	OrderID                   string `json:"order_id"`
	PayMethod                 int32  `json:"pay_method"`
}

func Charge(ctx context.Context, c *app.RequestContext) {
	userId, ok := ctx.Value("userId").(uint32)
	if !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
		return
	}
	resp, err := client.Charge(ctx, float32(amount), req.OrderID, userId, req.PayMethod)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"transaction_id": resp.TransactionId,
		"pay_url":        resp.PayUrl,
	})
}
