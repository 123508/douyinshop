package main

import "github.com/123508/douyinshop/pkg/errorno"

var NoOrderIdError = &errorno.BasicMessageError{Code: 404, Message: "没有用户id"}

var SearchOrderError = &errorno.BasicMessageError{Code: 500, Message: "查询订单错误"}

var SearchOrderDetailsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单详细信息错误"}

var SearchOrderLogsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单日志信息错误"}

var CancelOrderError = &errorno.BasicMessageError{Code: 500, Message: "取消订单错误"}

var DeliveryOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单发货失败"}

var ReceiveOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单变为待收货失败"}

var RejectionOrderError = &errorno.BasicMessageError{Code: 500, Message: "拒单失败"}

var ConfirmOrderError = &errorno.BasicMessageError{Code: 500, Message: "确认订单失败"}

var UnableRejectionOrderError = &errorno.BasicMessageError{Code: 500, Message: "当前状态不允许拒单,请注意"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}
