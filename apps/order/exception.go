package main

import "github.com/123508/douyinshop/pkg/errorno"

var SubmitOrderError = &errorno.BasicMessageError{Code: 500, Message: "订单提交失败"}

var SearchOrderError = &errorno.BasicMessageError{Code: 404, Message: "查询订单错误"}

var SearchOrderDetailsError = &errorno.BasicMessageError{Code: 500, Message: "查询订单详细信息错误"}

var SearchOrderLogsError = &errorno.BasicMessageError{Code: 404, Message: "查询订单日志错误"}

var CancelOrderError = &errorno.BasicMessageError{Code: 500, Message: "取消订单失败"}

var UnableChangeStatusError = &errorno.BasicMessageError{Code: 400, Message: "更新状态失败,该状态不允许被更新"}

var NotPayReminderError = &errorno.BasicMessageError{Code: 400, Message: "没有支付订单,无法提醒发货"}

var DeliveredReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家已经发货,无需提醒"}

var CancelReminderError = &errorno.BasicMessageError{Code: 400, Message: "订单已经取消,无法提醒发货"}

var RefundingReminderError = &errorno.BasicMessageError{Code: 400, Message: "退款中,无法提醒发货"}

var RefundedReminderError = &errorno.BasicMessageError{Code: 400, Message: "已经退款,无法提醒发货"}

var RejectionReminderError = &errorno.BasicMessageError{Code: 400, Message: "商家拒绝发货,无法提醒"}

var StatusError = &errorno.BasicMessageError{Code: 400, Message: "不允许的行为"}

var CompletionOrderError = &errorno.BasicMessageError{Code: 400, Message: "无法确认收货"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}
