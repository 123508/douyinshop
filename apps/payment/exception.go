package main

import "github.com/123508/douyinshop/pkg/errorno"

var NoShopping = &errorno.BasicMessageError{Message: "没有该购物记录"}

var NotSupportWechatPay = &errorno.BasicMessageError{Code: 400, Message: "暂时不支持微信支付"}
