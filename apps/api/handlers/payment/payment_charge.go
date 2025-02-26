package payment

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/cloudwego/hertz/pkg/app"
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
	userId, _ := c.Get("userId")
	var req Request
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
	}
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
	}
	resp, err := client.Charge(ctx, float32(amount), req.OrderID, userId.(uint32), req.PayMethod)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"transaction_id": resp.TransactionId,
		"pay_url":        resp.PayUrl,
	})
}
