package main

import "github.com/123508/douyinshop/pkg/errorno"

var RepeatedShop = &errorno.BasicMessageError{Code: 404, Message: "你已注册店铺，请勿重复注册"}

var ShopNotFound = &errorno.BasicMessageError{Code: 404, Message: "无法找到店铺"}

var ProductLoss = &errorno.BasicMessageError{Code: 404, Message: "商品信息丢失"}

var FailFetchProductList = &errorno.BasicMessageError{Code: 404, Message: "无法获取商品列表"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}
