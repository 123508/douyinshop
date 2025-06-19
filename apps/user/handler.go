package main

import (
	"context"
	"errors"
	"github.com/123508/douyinshop/apps/user/repository"
	"github.com/123508/douyinshop/apps/user/service"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	UserService service.UserService
}

func NewUserServiceImpl(db *gorm.DB, rds *redis.Client) *UserServiceImpl {
	return &UserServiceImpl{
		UserService: service.NewService(repository.NewRepository(db), rds),
	}
}

// Register implements the UserServiceImpl interface.
// 用户注册接口
// 如果两个密码不同,则返回空
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {

	if req.Password != req.ConfirmPassword {
		util.LogError("密码和注册密码不一致", "Register", nil)
		return nil, service.PasswordNotEqual
	}

	//初始化用户属性
	user1 := &models.User{
		Email:  req.Email,
		Name:   req.Nickname,
		Phone:  req.Phone,
		Gender: req.Gender,
	}

	//初始化用户密码
	user2 := &models.UserLogin{
		Password: service.Encryption(req.Password),
	}

	registryUserId, err := s.UserService.RegistryUser(ctx, user1, user2)

	if err != nil || registryUserId == 0 {
		return nil, err
	}

	return &user.RegisterResp{UserId: registryUserId}, nil
}

// Login implements the UserServiceImpl interface.
// 用户登录接口
// 如果账号密码正确，返回user_id
// 如果账号或密码错误，user_id为0
// 如果用户被删除也返回user_id=0
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email or password is null")
	}

	loginId, err := s.UserService.Login(ctx, req.Email, service.Encryption(req.Password))

	if err != nil || loginId == 0 {
		return nil, err
	}

	return &user.LoginResp{UserId: loginId}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
// 获取用户信息接口
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {

	row, err := s.UserService.GetUserInfo(ctx, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	//如果查询到该用户,返回用户信息
	return &user.GetUserInfoResp{
		Email:    row.Email,
		Nickname: row.Name,
		Avatar:   row.Avatar,
		Phone:    row.Phone,
		Gender:   row.Gender,
	}, nil
}

// Logout implements the UserServiceImpl interface.
// 用户登出接口
// 将用户当前的token设置为
func (s *UserServiceImpl) Logout(ctx context.Context, req *user.LogoutReq) (resp *user.Empty, err error) {

	if err = s.UserService.Logout(ctx, req.Token, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}

// Update implements the UserServiceImpl interface.
// 用户更新接口
// 更新用户资料并
// 绑定整个更新为事务,如果出错就进行回滚
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.Empty, err error) {

	if req.TargetUserId == 0 {
		return nil, service.UpdateUserInfoError
	}

	u := &models.User{
		ID:     req.TargetUserId,
		Name:   req.Name,
		Phone:  req.Phone,
		Gender: req.Gender,
		Avatar: req.Avatar,
	}

	if err = s.UserService.Update(ctx, u, req.RequestUserId); err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}

// Delete implements the UserServiceImpl interface.
// 删除用户接口
func (s *UserServiceImpl) Delete(ctx context.Context, req *user.DeleteReq) (resp *user.Empty, err error) {

	if err = s.UserService.Delete(ctx, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
// 对外暴露的负责分发令牌的接口
func (s *UserServiceImpl) DeliverTokenByRPC(ctx context.Context, req *user.DeliverTokenReq) (resp *user.DeliveryResp, err error) {

	//生成后端token
	token, err := s.UserService.DeliverToken(ctx, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	return &user.DeliveryResp{Token: token}, nil
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
// 对外暴露的验证令牌接口
// 如果redis中标记该令牌无效,返回错误响应
// 如果令牌无效,返回错误响应
// 如果令牌存活时间小于等于阈值,刷新令牌并返回成功响应
// 如果令牌存活时间大于阈值,直接返回成功响应
// 注意每次需要使用响应去接收token
func (s *UserServiceImpl) VerifyTokenByRPC(ctx context.Context, req *user.VerifyTokenReq) (resp *user.VerifyResp, err error) {

	if req.Token == "" {
		return nil, service.NilToken
	}

	token, uid, err := s.UserService.VerifyToken(ctx, req.Token, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	return &user.VerifyResp{
		Res:    true,
		Token:  token,
		UserId: uid,
	}, err
}

// ListUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersReq) (resp *user.ListUsersResp, err error) {
	//请求页长有问题
	if req.Page < 1 || req.PageSize < 1 {
		return nil, service.BadPageOrPageSize
	}

	listWithCache, err := s.UserService.GetUserList(ctx, req.RequestUserId, req.Page, req.PageSize, req.Filter)

	userList := make([]*user.User, 0, len(listWithCache))

	for _, v := range listWithCache {
		u := &user.User{
			Email:  v.Email,
			Name:   v.Name,
			Avatar: v.Avatar,
			Phone:  v.Phone,
			Gender: v.Gender,
			Status: v.Status,
			Id:     v.ID,
		}

		userList = append(userList, u)
	}

	total := uint32(len(userList))

	return &user.ListUsersResp{
		Users: userList,
		Total: total,
	}, nil
}

// ChangePassword implements the UserServiceImpl interface.
func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *user.ChangePasswordReq) (resp *user.Empty, err error) {

	//查询传入条件错误
	if req.TargetUserId == 0 {
		return nil, service.UserNotExists
	}

	if err = s.UserService.ChangePassword(ctx, req.TargetUserId, req.RequestUserId, req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}

	return &user.Empty{}, nil

}

// ForgotPassword implements the UserServiceImpl interface.
// 这里需要发送验证码
func (s *UserServiceImpl) ForgotPassword(ctx context.Context, req *user.ForgotPasswordReq) (resp *user.ForgotPasswordResp, err error) {
	//请求邮箱不正确
	if req.Email == "" {
		return nil, service.UserNotExists
	}

	uid, err := s.UserService.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &user.ForgotPasswordResp{UserId: uid}, nil
}

// VerifySmsCode implements the UserServiceImpl interface.
// 这里判定验证码是否成功,成功就返回一个一次性的ResetToken
func (s *UserServiceImpl) VerifySmsCode(ctx context.Context, req *user.VerifySmsCodeReq) (resp *user.VerifySmsCodeResp, err error) {

	ResetToken, err := s.UserService.VerifySmsCode(ctx, req.VerifyCode, req.Type, req.TargetUserId, req.RequestUserId)

	if err != nil {
		return nil, err
	}

	return &user.VerifySmsCodeResp{
		ResetToken: ResetToken,
	}, nil
}

// ResetPassword implements the UserServiceImpl interface.
// 这里进行重置密码,重置完成后删除ResetToken
func (s *UserServiceImpl) ResetPassword(ctx context.Context, req *user.ResetPasswordReq) (resp *user.Empty, err error) {

	if err = s.UserService.ResetPassword(ctx, req.ResetToken, req.Email, req.NewPassword, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}

// BindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) BindEmail(ctx context.Context, req *user.BindEmailReq) (resp *user.Empty, err error) {

	if err = s.UserService.BindEmail(ctx, req.Email, req.ResetToken, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}

// UnbindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnbindEmail(ctx context.Context, req *user.UnbindEmailReq) (resp *user.Empty, err error) {
	//清理缓存
	err = s.UserService.UnBindEmail(ctx, req.ResetToken, req.TargetUserId, req.RequestUserId)
	if err != nil {
		return nil, err
	}

	return &user.Empty{}, nil
}

// FreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) FreezeUser(ctx context.Context, req *user.FreezeUserReq) (resp *user.Empty, err error) {
	if err = s.UserService.FreezeUser(ctx, req.ResetToken, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}

// UnfreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnfreezeUser(ctx context.Context, req *user.UnfreezeUserReq) (resp *user.Empty, err error) {
	if err = s.UserService.UnfreezeUser(ctx, req.ResetToken, req.TargetUserId, req.RequestUserId); err != nil {
		return nil, err
	}
	return &user.Empty{}, nil
}

// PreBindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) PreBindEmail(ctx context.Context, req *user.PreBindEmailReq) (resp *user.Empty, err error) {
	// TODO: Your code here...
	return
}

// PreUnbindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) PreUnbindEmail(ctx context.Context, req *user.PreUnbindEmailReq) (resp *user.Empty, err error) {
	// TODO: Your code here...
	return
}

// PreFreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) PreFreezeUser(ctx context.Context, req *user.PreFreezeUserReq) (resp *user.Empty, err error) {
	// TODO: Your code here...
	return
}

// PreUnfreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) PreUnfreezeUser(ctx context.Context, req *user.PreUnFreezeUserReq) (resp *user.Empty, err error) {
	// TODO: Your code here...
	return
}
