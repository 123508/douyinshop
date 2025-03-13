package order

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	"github.com/123508/douyinshop/kitex_gen/order/order_common"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

func DetailShop(ctx context.Context, c *app.RequestContext) {

	// 获取并解析 user_id 参数
	value, exists := c.Get("userId")
	_, ok := value.(uint32)
	if !exists || !ok {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "userId must be a number",
		})
		return
	}

	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "orderId 参数错误",
		})
		return
	}

	var orderDetails []order_common.OrderDetail
	if err = c.Bind(&orderDetails); err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{
			"error": "订单详情解析失败",
		})
		return
	}

	orderResp, err := client.ShopDetail(ctx, uint32(orderId), orderDetails)
	if err != nil {
		errorno.DealWithError(err, c)
		return
	}

	c.JSON(consts.StatusOK, utils.H{
		"order": orderResp,
	})

}
