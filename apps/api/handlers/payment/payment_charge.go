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
	value, exists := c.Get("userId")
	userId, ok := value.(uint32)
	if !exists || !ok {
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
	}
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": err.Error(),
		})
	}
	resp, err := client.Charge(ctx, float32(amount), req.OrderID, userId, req.PayMethod)
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

		return
	}
	c.JSON(consts.StatusOK, utils.H{
		"transaction_id": resp.TransactionId,
		"pay_url":        resp.PayUrl,
	})
}
