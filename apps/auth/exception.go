package main

import "github.com/123508/douyinshop/pkg/errorno"

var UserRoleSearchError = &errorno.BasicMessageError{Code: 500, Message: "查询用户对应的角色列表错误"}

var RoleListSearchError = &errorno.BasicMessageError{Code: 500, Message: "查询角色列表错误"}

var SearchSuperError = &errorno.BasicMessageError{Code: 500, Message: "查询超级管理员权限错误"}

var BadPageOrPageSize = &errorno.BasicMessageError{Code: 400, Message: "请求页数或页长错误"}

var SearchRolePermissionError = &errorno.BasicMessageError{Code: 500, Message: "查询角色对应的权限列表错误"}
