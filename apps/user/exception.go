package main

import "github.com/123508/douyinshop/pkg/errorno"

var InvalidToken = &errorno.BasicMessageError{Code: 401, Message: "token已过期,请重新登录"}

var ParseTokenError = &errorno.BasicMessageError{Code: 402, Message: "解析token失败"}

var SignFailError = &errorno.BasicMessageError{Code: 500, Message: "token签名失败"}

var RedisConnectionError = &errorno.BasicMessageError{Code: 404, Message: "Redis数据库连接异常"}

var UserNotExists = &errorno.BasicMessageError{Code: 401, Message: "用户不存在"}

var PasswordNotEqual = &errorno.BasicMessageError{Code: 400, Message: "密码不匹配,请重新输入"}

var ErrorUsernameOrPassword = &errorno.BasicMessageError{Code: 404, Message: "用户名或密码错误"}

var NilToken = &errorno.BasicMessageError{Code: 404, Message: "请求Token为空"}

var UpdatePasswordError = &errorno.BasicMessageError{Code: 500, Message: "更新用户密码失败,请联系管理员"}

var SearchMySQLError = &errorno.BasicMessageError{Code: 500, Message: "查询数据库异常"}

var UserStatusError = &errorno.BasicMessageError{Code: 400, Message: "用户状态异常,不允许被修改"}

var UpdateUserInfoError = &errorno.BasicMessageError{Code: 500, Message: "更新用户信息错误"}
