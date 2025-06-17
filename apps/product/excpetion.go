package main

import "github.com/123508/douyinshop/pkg/errorno"

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}
