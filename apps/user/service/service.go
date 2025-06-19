package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/123508/douyinshop/apps/user/repository"
	"github.com/123508/douyinshop/kitex_gen/user"
	"github.com/123508/douyinshop/pkg/component/condition"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

type UserService interface {
	GetRedis() *redis.Client
	RegistryUser(ctx context.Context, u *models.User, uLogin *models.UserLogin) (uint64, error)
	Login(ctx context.Context, email, password string) (uint64, error)
	Logout(ctx context.Context, token string, TargetUserId, RequestUserId uint64) error
	GetUserInfo(ctx context.Context, TargetUserId, RequestUserId uint64) (*models.User, error)
	Update(ctx context.Context, user *models.User, RequestUserId uint64) error
	Delete(ctx context.Context, TargetUserId, RequestUserId uint64) error
	DeliverToken(ctx context.Context, TargetUserId, RequestUserId uint64) (string, error)
	VerifyToken(ctx context.Context, token string, TargetUserId, RequestUserId uint64) (string, uint64, error)
	GetUserList(ctx context.Context, RequestUserId uint64, page, pageSize uint32, filters []*user.Filter) ([]models.User, error)
	ChangePassword(ctx context.Context, TargetUserId, RequestUserId uint64, OldPwd, NewPwd string) error
	ForgotPassword(ctx context.Context, Email string) (uint64, error)
	VerifySmsCode(ctx context.Context, SmsCode string, Type, TargetUserId, RequestUserId uint64) (string, error)
	ResetPassword(ctx context.Context, ResetToken, email, NewPwd string, TargetUserId, RequestUserId uint64) error
	BindEmail(ctx context.Context, Email, ResetToken string, TargetUserId, RequestUserId uint64) error
	UnBindEmail(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error
	FreezeUser(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error
	UnfreezeUser(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error
	PreBindEmail(ctx context.Context, TargetUserId, RequestUserId uint64) error
	PreUnbindEmail(ctx context.Context, TargetUserId, RequestUserId uint64) error
	PreFreezeUser(ctx context.Context, TargetUserId, RequestUserId uint64) error
	PreUnfreezeUser(ctx context.Context, TargetUserId, RequestUserId uint64) error
}

type ServiceImpl struct {
	userRepo repository.UserRepository
	Rds      *redis.Client
}

func NewService(userRepo repository.UserRepository, rds *redis.Client) UserService {
	return &ServiceImpl{
		userRepo: userRepo,
		Rds:      rds,
	}
}

func (s *ServiceImpl) GetRedis() *redis.Client {
	return s.Rds
}

func (s *ServiceImpl) RegistryUser(ctx context.Context, u *models.User, uLogin *models.UserLogin) (uint64, error) {
	createUserId, err := s.userRepo.CreateUser(ctx, u, uLogin)
	if err != nil {
		return 0, RegisterUserError
	}
	return createUserId, nil
}

func (s *ServiceImpl) Login(ctx context.Context, email, password string) (uint64, error) {

	infoFromEmail, err := s.userRepo.GetUserInfoFromEmail(ctx, email)

	if err != nil {
		return 0, ErrorUsernameOrPassword
	}

	ok, err := s.userRepo.CompareUserPassword(ctx, infoFromEmail.ID, password)

	if err != nil || !ok {
		return 0, ErrorUsernameOrPassword
	}

	return infoFromEmail.ID, nil
}

func (s *ServiceImpl) GetUserInfo(ctx context.Context, TargetUserId, RequestUserId uint64) (*models.User, error) {

	simple := util.SimpleCacheComponent[uint64, *models.User]{
		Rds:       s.Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(serviceName, TargetUserId),
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		FuncName:  "GetUserInfo",
		QueryExec: func() (*models.User, error) {
			return s.userRepo.GetUserInfoFromId(ctx, TargetUserId)
		},
		Expires: time.Duration(rand.Intn(10)+5) * time.Minute,
	}

	row, err := simple.QueryWithCache()

	if err != nil {
		return nil, err
	}

	return row, nil
}

func (s *ServiceImpl) Logout(ctx context.Context, token string, TargetUserId, RequestUserId uint64) error {
	//让token失效,否则报错并返回
	if err := s.Rds.Set(ctx, token, "1", 8*time.Hour).Err(); err != nil {
		return TokenTimeOutError
	}

	if err := s.Rds.Del(ctx, util.TakeKey("SToken", TargetUserId)).Err(); err != nil {
		return TokenTimeOutError
	}

	return nil
}

func (s *ServiceImpl) Update(ctx context.Context, user *models.User, RequestUserId uint64) error {

	if user == nil || user.ID == 0 {
		return UserNotExists
	}

	//清理缓存
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, user.ID))

	info, err := s.userRepo.GetUserInfoFromId(ctx, user.ID)

	if err != nil {
		return err
	}

	if info.Status == 1 {
		return UserFreezeError
	}

	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return UpdateUserInfoError
	}

	return nil
}

func (s *ServiceImpl) Delete(ctx context.Context, TargetUserId, RequestUserId uint64) error {
	//清理缓存
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, TargetUserId))

	if err := s.userRepo.DeleteUser(ctx, TargetUserId); err != nil {
		return DeleteUserError
	}

	return nil
}

func (s *ServiceImpl) DeliverToken(ctx context.Context, TargetUserId, RequestUserId uint64) (string, error) {
	//生成后端token
	token, err := GenerateFrontendJWT(TargetUserId)
	if err != nil {
		return "", err
	}

	//这里需要调用auth服务去获取必要的信息
	roleCodes := make([]string, 0)
	permCodes := make([]string, 0)

	//生成并存储后端token
	backendJWT, err := GenerateBackendJWT(TargetUserId, roleCodes, permCodes, 0)

	if err != nil {
		return "", err
	}

	//过期时间定为7天
	s.Rds.Set(ctx, util.TakeKey("SToken", TargetUserId), backendJWT, 7*24*time.Hour)

	return token, nil
}

func (s *ServiceImpl) VerifyToken(ctx context.Context, token string, TargetUserId, RequestUserId uint64) (string, uint64, error) {
	if token == "" {
		return "", 0, NilToken
	}

	//如果这个token已经进入了黑名单,就直接返回错误
	isExist, err := s.Rds.Exists(ctx, token).Result()

	if err != nil || isExist == 1 {
		return "", 0, TokenTimeOutError
	}

	//解析前端token
	FrontendJwt, err := ParseFrontendJWT(token)

	//解析前端token失败,返回异常
	if err != nil {
		return "", 0, err
	}

	key := util.TakeKey("SToken", FrontendJwt.UserId)

	//在redis中检查token是否存在
	ServerToken, err := s.Rds.Get(ctx, key).Result()

	//如果出错,或者token不存在,直接返回错误信息
	if err != nil || ServerToken == "" {
		util.LogError("Token不存在", "VerifyTokenByRPC", err)
		return "", 0, TokenTimeOutError
	}

	//判断令牌是否可以被解析,如果令牌无法被解析返回失败响应
	ServerJwt, err := ParseBackendJWT(ServerToken)

	if err != nil {
		return "", 0, err
	}

	//签发人错误,有可能是伪造令牌
	if ServerJwt.Issuer != FrontendJwt.Issuer {
		return "", 0, IssuerNotMatchError
	}

	var respToken string

	//如果相差时间小于令牌存活阈值,就重新生成前端令牌
	suv := time.Duration(max[int](config.Conf.AdminTtl-config.Conf.AdminSuv, 10800)) * time.Second

	if time.Since(ServerJwt.IssuedAt.Time) >= suv {
		newToken, err := GenerateFrontendJWT(ServerJwt.UserId)
		if err != nil {
			return "", 0, err // 返回错误
		}
		//将token重新放入
		respToken = newToken

		s.Rds.Set(ctx, token, "1", 7*24*time.Hour)
	} else {
		respToken = token
	}

	return respToken, ServerJwt.UserId, err
}

func (s *ServiceImpl) GetUserList(ctx context.Context, RequestUserId uint64,
	page, pageSize uint32, filters []*user.Filter) ([]models.User, error) {
	//请求页长有问题
	if page < 1 || pageSize < 1 {
		return nil, BadPageOrPageSize
	}

	builder := condition.NewConditionBuilder()

	for _, v := range filters {
		builder = builder.And(v.FieldName, v.Operator, v.Value)
	}

	sql, params := builder.Build().ToSQL()

	listQuery := util.ListCacheComponent[uint64, models.User]{
		Rds:             s.Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(serviceName, util.Md5Hash(util.TakeKey(pageSize, pageSize))),
		DetailKeyPrefix: util.TakeKey(serviceName),
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FuncName:        "ListUsers",
		FullQueryExec: func() ([]models.User, error) {
			return s.userRepo.GetUserList(ctx, int(page), int(pageSize), sql, params)
		},
		Expires:     time.Duration(rand.Intn(3)+3) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          s.userRepo.GetDB(),
	}

	listWithCache, err := listQuery.QueryListWithCache()

	if err != nil {
		return nil, GetUserListError
	}

	return listWithCache, nil
}

func (s *ServiceImpl) ChangePassword(ctx context.Context, TargetUserId, RequestUserId uint64, OldPwd, NewPwd string) error {
	//查询传入条件错误
	if TargetUserId == 0 {
		return UserNotExists
	}

	//清理缓存
	defer util.CleanCache(s.Rds, ctx, util.TakeKey(serviceName, TargetUserId))

	//用户状态合法性校验
	info, err := s.IsLegalUserById(ctx, TargetUserId)

	if err != nil || info == nil {
		return err
	}

	//查询用户密码
	ok, err := s.userRepo.CompareUserPassword(ctx, TargetUserId, Encryption(OldPwd))

	//用户密码不正确
	if err != nil || !ok {
		return PasswordNotEqual
	}

	if err = s.userRepo.UpdateUserPassword(ctx, TargetUserId, Encryption(NewPwd)); err != nil {
		return UpdatePasswordError
	}

	return nil
}

func (s *ServiceImpl) ForgotPassword(ctx context.Context, email string) (uint64, error) {
	//请求邮箱不正确
	if email == "" {
		return 0, UserNotExists
	}

	//用户状态合法性校验
	info, err := s.IsLegalUserByEmail(ctx, email)

	if err != nil || info == nil {
		return 0, SearchUserError
	}

	vCode := rand.Intn(6)

	if err = s.Rds.Set(ctx, util.TakeKey(serviceName, "Sms", info.ID, util.Md5Hash("0")), vCode, 10*time.Minute).Err(); err != nil {
		return 0, RedisSetError
	}

	s.SendSmsCode(ctx, vCode)

	fmt.Println("发送验证码成功:", vCode, "用户邮箱为:", email, "操作类型:忘记密码")

	return info.ID, nil
}

func (s *ServiceImpl) VerifySmsCode(ctx context.Context, SmsCode string, Type, TargetUserId, RequestUserId uint64) (string, error) {
	//获取存储好的验证码 Type: 0 密码  1绑定邮箱  2解绑邮箱 3冻结用户 4解冻用户
	vCodeKey := util.TakeKey(serviceName, "Sms", TargetUserId, util.Md5Hash(strconv.FormatUint(Type, 10)))
	vCode, err := s.Rds.Get(ctx, vCodeKey).Result()

	//redis查询失败
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", SearchRedisError
	}

	//验证码已过期的情况
	if vCode == "" {
		return "", VerifyCodeTimeOutError
	}

	info, err := s.IsLegalUserById(ctx, TargetUserId)

	if err != nil || info == nil {
		return "", err
	}

	//验证码错误
	if vCode != SmsCode {
		return "", VerifyCodeError
	}

	v7, _ := uuid.NewV7()

	ResetToken := v7.String()

	//写入重置token
	if err = s.Rds.Set(ctx, ResetToken, "1", 10*time.Minute).Err(); err != nil {
		return "", RedisSetError
	}

	//删除对应的短信验证码
	s.Rds.Del(ctx, vCodeKey)

	return ResetToken, nil
}

func (s *ServiceImpl) ResetPassword(ctx context.Context, ResetToken, email, NewPwd string, TargetUserId, RequestUserId uint64) error {
	//获取存储好的验证码
	_, err := s.Rds.Get(ctx, ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return VerifyCodeTimeOutError
	}

	//redis查询失败
	if err != nil {
		return SearchRedisError
	}

	info, err := s.IsLegalUserByEmail(ctx, email)

	if err != nil || info == nil {
		return SearchUserError
	}

	if err = s.userRepo.UpdateUserPassword(ctx, TargetUserId, Encryption(NewPwd)); err != nil {
		return UpdatePasswordError
	}

	//删除重置Token,保证其只可用一次
	s.Rds.Del(ctx, ResetToken)

	return nil
}

func (s *ServiceImpl) PreBindEmail(ctx context.Context, TargetUserId, RequestUserId uint64) error {
	if TargetUserId == 0 {
		return UserNotExists
	}

	info, err := s.IsLegalUserById(ctx, TargetUserId)

	if err != nil || info == nil {
		return SearchUserError
	}

	vCode := rand.Intn(6)

	if err = s.Rds.Set(ctx, util.TakeKey(serviceName, "Sms", TargetUserId, util.Md5Hash("1")), vCode, 10*time.Minute).Err(); err != nil {
		return RedisSetError
	}

	s.SendSmsCode(ctx, vCode)

	fmt.Println("发送验证码成功:", vCode, "用户Id为:", TargetUserId, "操作类型:绑定邮箱")

	return nil
}

func (s *ServiceImpl) BindEmail(ctx context.Context, Email, ResetToken string, TargetUserId, RequestUserId uint64) error {
	//获取存储好的验证码
	_, err := s.Rds.Get(ctx, ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return VerifyCodeTimeOutError
	}

	defer s.Rds.Del(ctx, util.TakeKey(serviceName, TargetUserId))

	info, err := s.IsLegalUserById(ctx, TargetUserId)

	if err != nil {
		return err
	}

	//存在邮箱则不允许直接修改
	if info.Email != "" {
		return NotAllowedBindEmail
	}

	//绑定邮箱
	if err = s.userRepo.SetEmail(ctx, TargetUserId, Email); err != nil {
		return BindEmailError
	}

	//删除Token,保证其只可用一次
	s.Rds.Del(ctx, ResetToken)

	return nil
}

func (s *ServiceImpl) PreUnbindEmail(ctx context.Context, TargetUserId, RequestUserId uint64) error {
	if TargetUserId == 0 {
		return UserNotExists
	}

	info, err := s.IsLegalUserById(ctx, TargetUserId)

	if err != nil || info == nil {
		return SearchUserError
	}

	vCode := rand.Intn(6)

	if err = s.Rds.Set(ctx, util.TakeKey(serviceName, "Sms", TargetUserId, util.Md5Hash("2")), vCode, 10*time.Minute).Err(); err != nil {
		return RedisSetError
	}

	s.SendSmsCode(ctx, vCode)

	fmt.Println("发送验证码成功:", vCode, "用户Id为:", TargetUserId, "操作类型:解绑邮箱")

	return nil
}

func (s *ServiceImpl) UnBindEmail(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error {

	//获取存储好的验证码
	_, err := s.Rds.Get(ctx, ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return VerifyCodeTimeOutError
	}

	//清理缓存
	defer s.Rds.Del(ctx, util.TakeKey(serviceName, TargetUserId))

	if _, err = s.IsLegalUserById(ctx, TargetUserId); err != nil {
		return err
	}

	if err = s.userRepo.RemoveEmail(ctx, TargetUserId); err != nil {
		return UnBindEmailError
	}

	//删除Token,保证其只可用一次
	s.Rds.Del(ctx, ResetToken)

	return nil
}

func (s *ServiceImpl) PreFreezeUser(ctx context.Context, TargetUserId, RequestUserId uint64) error {

	if TargetUserId == 0 {
		return UserNotExists
	}

	info, err := s.userRepo.GetUserInfoFromId(ctx, TargetUserId)

	if err != nil || info == nil {
		return SearchUserError
	}

	vCode := rand.Intn(6)

	if err = s.Rds.Set(ctx, util.TakeKey(serviceName, "Sms", TargetUserId, util.Md5Hash("3")), vCode, 10*time.Minute).Err(); err != nil {
		return RedisSetError
	}

	s.SendSmsCode(ctx, vCode)

	fmt.Println("发送验证码成功:", vCode, "用户Id为:", TargetUserId, "操作类型:冻结用户")

	return nil
}

func (s *ServiceImpl) FreezeUser(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error {
	//获取存储好的验证码
	_, err := s.Rds.Get(ctx, ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return VerifyCodeTimeOutError
	}

	//清理缓存
	defer s.Rds.Del(ctx, util.TakeKey(serviceName, TargetUserId))

	info, err := s.userRepo.GetUserInfoFromId(ctx, TargetUserId)

	if err != nil {
		return err
	}

	//状态已经修改
	if info.Status == 1 {
		return nil
	}

	if err = s.userRepo.FreezeUser(ctx, TargetUserId); err != nil {
		return UpdateUserInfoError
	}

	//删除Token,保证其只可用一次
	s.Rds.Del(ctx, ResetToken)

	return nil
}

func (s *ServiceImpl) PreUnfreezeUser(ctx context.Context, TargetUserId, RequestUserId uint64) error {
	if TargetUserId == 0 {
		return UserNotExists
	}

	info, err := s.userRepo.GetUserInfoFromId(ctx, TargetUserId)

	if err != nil || info == nil {
		return SearchUserError
	}

	vCode := rand.Intn(6)

	if err = s.Rds.Set(ctx, util.TakeKey(serviceName, "Sms", TargetUserId, util.Md5Hash("4")), vCode, 10*time.Minute).Err(); err != nil {
		return RedisSetError
	}

	s.SendSmsCode(ctx, vCode)

	fmt.Println("发送验证码成功:", vCode, "用户Id为:", TargetUserId, "操作类型:解冻用户")

	return nil
}

func (s *ServiceImpl) UnfreezeUser(ctx context.Context, ResetToken string, TargetUserId, RequestUserId uint64) error {

	//获取存储好的验证码
	_, err := s.Rds.Get(ctx, ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return VerifyCodeTimeOutError
	}
	//清理缓存
	defer s.Rds.Del(ctx, util.TakeKey(serviceName, TargetUserId))

	info, err := s.userRepo.GetUserInfoFromId(ctx, TargetUserId)

	if err != nil {
		return err
	}

	//状态已经修改
	if info.Status == 0 {
		return nil
	}

	if err = s.userRepo.UnfreezeUser(ctx, TargetUserId); err != nil {
		return UpdateUserInfoError
	}

	//删除Token,保证其只可用一次
	s.Rds.Del(ctx, ResetToken)
	return nil
}

func (s *ServiceImpl) IsLegalUserById(ctx context.Context, TargetUserId uint64) (*models.User, error) {
	//用户状态合法性校验
	info, err := s.userRepo.GetUserInfoFromId(ctx, TargetUserId)

	//查询错误
	if err != nil {
		return nil, SearchUserError
	}

	//用户被冻结
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	return info, nil
}

func (s *ServiceImpl) IsLegalUserByEmail(ctx context.Context, Email string) (*models.User, error) {
	info, err := s.userRepo.GetUserInfoFromEmail(ctx, Email)

	if err != nil {
		return nil, SearchUserError
	}

	if info.Status == 1 {
		return nil, UserFreezeError
	}

	return info, nil
}

func (s *ServiceImpl) SendSmsCode(ctx context.Context, vCode any) {
	//TODO 这里之后调用发送验证码的逻辑

	fmt.Println(vCode)
}
