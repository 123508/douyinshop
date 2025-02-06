package main

import (
	"context"
	"strconv"

	user "github.com/123508/douyinshop/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	// 临时测试用的代码，实际开发中需要替换为真实的逻辑
	id, err := strconv.ParseInt(req.Email, 10, 32)
	if err != nil {
		id = 0
	}
	resp = &user.LoginResp{
		UserId: int32(id),
	}

	return
}
