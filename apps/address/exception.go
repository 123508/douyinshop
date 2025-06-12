package main

import "github.com/123508/douyinshop/pkg/errorno"

var InvalidAddressIdError = &errorno.BasicMessageError{Code: 400, Message: "地址ID无效"}

var ForbiddenDeleteError = &errorno.BasicMessageError{Code: 401, Message: "地址不存在或无权限删除"}

var DeleteAddrError = &errorno.BasicMessageError{Code: 500, Message: "地址删除失败,请联系管理员"}

var ForbiddenAskError = &errorno.BasicMessageError{Code: 401, Message: "无权访问该地址"}

var FailUpdateError = &errorno.BasicMessageError{Code: 500, Message: "更新地址失败,请联系管理员"}

var GetDefaultError = &errorno.BasicMessageError{Code: 404, Message: "获取默认地址失败"}
