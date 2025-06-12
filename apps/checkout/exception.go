package main

import "github.com/123508/douyinshop/pkg/errorno"

var GetCartError = &errorno.BasicMessageError{Code: 404, Message: "获取购物车错误"}

var CartNotExistError = &errorno.BasicMessageError{Code: 404, Message: "购物车不存在"}

var SubmitError = &errorno.BasicMessageError{Code: 500, Message: "提交错误"}

var ChargeError = &errorno.BasicMessageError{Code: 417, Message: "支付异常"}
