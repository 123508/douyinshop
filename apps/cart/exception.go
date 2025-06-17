package main

import "github.com/123508/douyinshop/pkg/errorno"

var NilRequestError = &errorno.BasicMessageError{Code: 400, Message: "请求为空"}

var NilItemError = &errorno.BasicMessageError{Code: 400, Message: "商品为空"}

var NegativeQuantityError = &errorno.BasicMessageError{Code: 400, Message: "数量必须为正数"}

var NilUserIdError = &errorno.BasicMessageError{Code: 400, Message: "用户id不能为空"}

var FailToUpdateCart = &errorno.BasicMessageError{Code: 403, Message: "无法更新购物车"}

var NilProductIdError = &errorno.BasicMessageError{Code: 400, Message: "商品id不能为空"}
