package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/123508/douyinshop/pkg/component/condition"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/123508/douyinshop/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
// 用户注册接口
// 如果两个密码不同,则返回空
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	if req.Password != req.ConfirmPassword {
		util.LogError("密码和注册密码不一致", "Register", nil)
		return nil, PasswordNotEqual
	}

	//初始化用户属性
	user1 := &models.User{}
	user1.Email = req.Email
	user1.Name = req.Nickname
	user1.Phone = req.Phone
	user1.Gender = req.Gender
	user2 := &models.UserLogin{}
	user2.Password = Encryption(req.Password)

	err = DB.Transaction(func(tx *gorm.DB) error {

		//创建用户对象
		if err := DB.Create(user1).Error; err != nil {
			return err
		}

		user2.UserId = user1.ID

		if err := DB.Create(user2).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		util.LogError("创建用户对象错误", "Register", err)
		return nil, err
	}

	return &user.RegisterResp{UserId: user1.ID}, nil
}

// Login implements the UserServiceImpl interface.
// 用户登录接口
// 如果账号密码正确，返回user_id
// 如果账号或密码错误，user_id为0
// 如果用户被删除也返回user_id=0
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	//从users表中读取Email信息
	var row models.User
	DB.Model(&models.User{}).Where("email = ?", req.Email).Find(&row)

	//从user_logins表中获取Password信息
	var res models.UserLogin
	DB.Model(&models.UserLogin{}).
		Where("user_id = ? and password = ?", row.ID, Encryption(req.Password)).
		Find(&res)

	//如果用户已经被删除也返回空
	if res.ID == 0 {
		return &user.LoginResp{UserId: 0}, ErrorUsernameOrPassword
	}

	//这个方法之后要调用发放token的逻辑

	return &user.LoginResp{UserId: row.ID}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
// 获取用户信息接口
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {

	simple := util.SimpleCacheComponent[uint64, models.User]{
		Rds:       Rds,
		Ctx:       ctx,
		Key:       util.TakeKey(serviceName, req.UserId),
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		FuncName:  "GetUserInfo",
		QueryExec: func() (models.User, error) {
			var row models.User
			if err := DB.Model(&models.User{}).
				Where("id = ?", req.UserId).
				First(&row).Error; err != nil {
				util.LogError("用户不存在", "GetUserInfo", err)
				return models.User{}, UserNotExists
			}
			return row, nil
		},
		Expires: time.Duration(rand.Intn(10)+5) * time.Minute,
	}

	row, err := simple.QueryWithCache()

	if err != nil {
		return nil, err
	}

	////通用类型
	//var row models.User
	//
	////查询缓存
	//result, _ := Rds.Get(ctx, util.TakeKey(serviceName, req.UserId)).Result()
	//
	//err = json.Unmarshal([]byte(result), &row)
	//
	////查询缓存失败
	//if err != nil {
	//
	//	//如果查询到不存在该用户,返回error
	//	if err := DB.Model(&models.User{}).Where("id = ?", req.UserId).First(&row).Error; err != nil {
	//		util.LogError("用户不存在", "GetUserInfo", "", err)
	//		return nil, UserNotExists
	//	}
	//
	//	//存入缓存
	//	jsonData, _ := json.Marshal(row)
	//	err := Rds.Set(ctx, util.TakeKey(serviceName, req.UserId), string(jsonData), time.Duration(rand.Intn(15)+30)*time.Minute).Err()
	//	if err != nil {
	//		util.LogError("存入缓存失败", "GetUserInfo", "", err)
	//	}
	//} else {
	//	log.WithFields(log.Fields{
	//		"方法名": "GetUserInfo",
	//	}).Info("查询缓存成功")
	//}

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

	//让token失效,否则报错并返回
	if err = Rds.Set(ctx, req.Token, "1", 8*time.Hour).Err(); err != nil {
		return &user.Empty{}, err
	}

	return &user.Empty{}, nil

}

// Update implements the UserServiceImpl interface.
// 用户更新接口
// 更新用户资料并
// 绑定整个更新为事务,如果出错就进行回滚
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.Empty, err error) {

	if req.UserId == 0 {
		return nil, UserNotExists
	}

	//清理缓存
	defer util.CleanCache(Rds, ctx, util.TakeKey(serviceName, req.UserId))

	info, err := GetUserInfo(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//更新用户信息
	updates := make(map[string]interface{}, 5)

	if req.Gender != info.Gender {
		updates["gender"] = req.Gender
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Nickname != "" {
		updates["name"] = req.Nickname
	}
	if err := DB.Model(&models.User{}).
		Where("id=?", req.UserId).
		Updates(updates).Error; err != nil {
		util.LogError("更新用户出错", "Update", err)
		return nil, UpdateUserInfoError
	}

	return &user.Empty{}, nil
}

// Delete implements the UserServiceImpl interface.
// 删除用户接口
// 绑定事务
func (s *UserServiceImpl) Delete(ctx context.Context, req *user.DeleteReq) (resp *user.Empty, err error) {

	//清理缓存
	defer util.CleanCache(Rds, ctx, util.TakeKey(serviceName, req.UserId))

	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := DB.Model(&models.User{}).Where("id = ?", req.UserId).Update("phone", nil).Update("email", nil).Delete(&models.User{}).Error; err != nil {
			return err
		}
		if err := DB.Where("user_id=?", req.UserId).Delete(&models.UserLogin{}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		klog.Error("删除用户异常")
		return nil, err
	}
	return &user.Empty{}, nil
}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
// 对外暴露的负责分发令牌的接口
func (s *UserServiceImpl) DeliverTokenByRPC(ctx context.Context, req *user.DeliverTokenReq) (resp *user.DeliveryResp, err error) {

	//生成后端token
	token, err := GenerateFrontendJWT(req.UserId)
	if err != nil {
		return nil, err
	}

	//生成并存储后端token
	backendJWT, err := GenerateBackendJWT(req.UserId, req.RoleCodes, req.PermCodes, 0)

	if err != nil {
		return nil, err
	}

	//过期时间定为7天
	Rds.Set(ctx, util.TakeKey("SToken", req.UserId), backendJWT, 7*24*time.Hour)

	resp = &user.DeliveryResp{Token: token}
	return resp, nil
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
		return nil, NilToken
	}

	//如果这个token已经进入了黑名单,就直接返回
	isExist, err := Rds.Exists(ctx, req.Token).Result()

	if err != nil || isExist == 1 {
		return nil, InvalidToken
	}

	//解析前端token
	FrontendJwt, err := ParseFrontendJWT(req.Token)

	//解析前端token失败,返回异常
	if err != nil {
		return nil, err
	}

	key := util.TakeKey("SToken", FrontendJwt.UserId)

	//在redis中检查token是否存在
	ServerToken, err := Rds.Get(ctx, key).Result()

	//如果出错,或者token不存在,直接返回错误信息
	if err != nil || ServerToken == "" {
		util.LogError("Token不存在", "VerifyTokenByRPC", err)
		return nil, InvalidToken
	}

	//判断令牌是否可以被解析,如果令牌无法被解析返回失败响应
	ServerJwt, err := ParseBackendJWT(ServerToken)

	if err != nil {
		return nil, err
	}

	//签发人错误,有可能是伪造令牌
	if ServerJwt.Issuer != FrontendJwt.Issuer {
		return nil, IssuerNotMatch
	}

	//如果相差时间小于令牌存活阈值,就重新生成前端令牌
	suv := time.Duration(max[int](config.Conf.AdminTtl-config.Conf.AdminSuv, 10800)) * time.Second

	resp = &user.VerifyResp{}

	if time.Since(ServerJwt.IssuedAt.Time) >= suv {
		newToken, err := GenerateFrontendJWT(ServerJwt.UserId)
		if err != nil {
			return nil, err // 返回错误
		}
		//将token重新放入
		resp.Token = newToken

		Rds.Set(ctx, req.Token, "1", 7*24*time.Hour)
	} else {
		resp.Token = req.Token
	}
	resp.UserId = ServerJwt.UserId

	return resp, err
}

// ListUsers implements the UserServiceImpl interface.
func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersReq) (resp *user.ListUsersResp, err error) {
	//请求页长有问题
	if req.Page < 1 || req.PageSize < 1 {
		return nil, BadPageOrPageSize
	}

	builder := condition.NewConditionBuilder()

	for _, v := range req.Filter {
		builder = builder.And(v.FieldName, v.Operator, v.Value)
	}

	sql, params := builder.Build().ToSQL()

	listQuery := util.ListCacheComponent[uint64, models.User]{
		Rds:             Rds,
		Ctx:             ctx,
		IdListKey:       util.TakeKey(serviceName, util.Md5Hash(util.TakeKey(req.PageSize, req.PageSize))),
		DetailKeyPrefix: util.TakeKey(serviceName),
		Marshal:         json.Marshal,
		Unmarshal:       json.Unmarshal,
		FuncName:        "ListUsers",
		FullQueryExec: func() ([]models.User, error) {
			var users []models.User
			offset := int((req.Page - 1) * req.PageSize)
			if err := DB.Model(&models.User{}).
				Where(sql, params...).
				Offset(offset).
				Limit(int(req.PageSize)).
				Find(&users).Error; err != nil {
				return nil, err
			}
			return users, nil
		},
		Expires:     time.Duration(rand.Intn(3)+3) * time.Minute,
		MaxLostRate: 30,
		Sort:        nil,
		IdName:      "id",
		DB:          DB,
	}

	listWithCache, err := listQuery.QueryListWithCache()

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
	if req.UserId == 0 {
		return nil, UserNotExists
	}

	//清理缓存
	defer util.CleanCache(Rds, ctx, util.TakeKey(serviceName, req.UserId))

	//用户状态合法性校验
	info, err := GetUserInfo(ctx, req.UserId)

	//查询错误
	if err != nil {
		return nil, err
	}

	//用户被冻结
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//查询用户密码
	var userLogin models.UserLogin

	if err := DB.Model(&models.UserLogin{}).Where("user_id = ?", req.UserId).First(&userLogin).Error; err != nil {
		util.LogError("查询用户密码异常", "GetUserPassword", err)
		return nil, err
	}

	//用户密码不正确
	if userLogin.Password != Encryption(req.OldPassword) {
		return nil, PasswordNotEqual
	}

	//修改密码部分
	where := DB.Model(&models.UserLogin{}).Where("user_id = ?", req.UserId)
	if err = where.Update("password", Encryption(req.NewPassword)).Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("更新用户密码错误", "ChangePassword", err)
		return nil, UpdatePasswordError
	}

	return &user.Empty{}, nil
}

// ForgotPassword implements the UserServiceImpl interface.
// 这里需要发送验证码
func (s *UserServiceImpl) ForgotPassword(ctx context.Context, req *user.ForgotPasswordReq) (resp *user.Empty, err error) {
	//请求邮箱不正确
	if req.Email == "" {
		return nil, UserNotExists
	}

	//用户状态合法性校验
	var info models.User

	if err = DB.Model(&models.User{}).Where("email = ?", req.Email).Find(&info).Error; err != nil {
		return nil, err
	}

	//用户状态异常
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	vCode := rand.Intn(6)

	if err = Rds.Set(ctx, util.TakeKey(serviceName, req.Email, util.Md5Hash("params")), vCode, 10*time.Minute).Err(); err != nil {
		return nil, RedisSetError
	} else {
		fmt.Println("发送验证码成功:", vCode, "用户邮箱为:", req.Email)
	}

	//TODO 这里之后调用发送验证码的逻辑

	return &user.Empty{}, nil
}

// VerifySmsCode implements the UserServiceImpl interface.
// 这里判定验证码是否成功,成功就返回一个一次性的ResetToken
func (s *UserServiceImpl) VerifySmsCode(ctx context.Context, req *user.VerifySmsCodeReq) (resp *user.VerifySmsCodeResp, err error) {

	//获取存储好的验证码
	vCodeKey := util.TakeKey(serviceName, req.Email, util.Md5Hash("params"))
	vCode, err := Rds.Get(ctx, vCodeKey).Result()

	//redis查询失败
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, SearchRedisError
	}

	//验证码已过期的情况
	if vCode == "" {
		return nil, VerifyCodeTimeOutError
	}

	//用户状态合法性校验
	var info models.User

	if err = DB.Model(&models.User{}).Where("email = ?", req.Email).Find(&info).Error; err != nil {
		return nil, err
	}

	//用户状态异常
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//验证码错误
	if vCode != req.VerifyCode {
		return nil, VerifyCodeError
	}

	ResetToken := util.Md5Hash(util.TakeKey(serviceName, req.Email, util.Md5Hash("params")))

	//写入重置token
	if err = Rds.Set(ctx, ResetToken, "1", 10*time.Minute).Err(); err != nil {
		return nil, RedisSetError
	}

	//删除对应的短信验证码
	Rds.Del(ctx, vCodeKey)

	return &user.VerifySmsCodeResp{
		ResetToken: ResetToken,
	}, nil
}

// ResetPassword implements the UserServiceImpl interface.
// 这里进行重置密码,重置完成后删除ResetToken
func (s *UserServiceImpl) ResetPassword(ctx context.Context, req *user.ResetPasswordReq) (resp *user.Empty, err error) {

	//获取存储好的验证码
	_, err = Rds.Get(ctx, req.ResetToken).Result()

	//验证码已过期的情况
	if errors.Is(err, redis.Nil) {
		return nil, VerifyCodeTimeOutError
	}

	//redis查询失败
	if err != nil {
		return nil, SearchRedisError
	}

	//用户状态合法性校验
	var info models.User

	if err = DB.Model(&models.User{}).Where("email = ?", req.Email).Find(&info).Error; err != nil {
		return nil, err
	}

	//用户状态异常
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//更新用户密码
	where := DB.Model(&models.UserLogin{}).Where("id = ?", info.ID)
	if err = where.Update("password", Encryption(req.NewPassword)).Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("更新用户密码错误", "ChangePassword", err)
		return nil, UpdatePasswordError
	}

	//删除重置Token,保证其只可用一次
	Rds.Del(ctx, req.ResetToken)

	return &user.Empty{}, nil
}

// BindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) BindEmail(ctx context.Context, req *user.BindEmailReq) (resp *user.Empty, err error) {

	defer Rds.Del(ctx, util.TakeKey(serviceName, req.UserId))

	info, err := GetUserInfo(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	//用户为冻结态
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//存在邮箱则不允许直接修改
	if info.Email != "" {
		return nil, NotAllowedBindEmail
	}

	//绑定邮箱
	if err := DB.Model(&models.User{}).Where("id = ?", req.UserId).
		Update("updated_at", time.Now()).
		Update("email", req.Email).Error; err != nil {
		util.LogError("绑定邮箱错误", "BindEmail", err)
		return nil, err
	}

	return
}

// UnbindEmail implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnbindEmail(ctx context.Context, req *user.UnbindEmailReq) (resp *user.Empty, err error) {
	//清理缓存
	defer Rds.Del(ctx, util.TakeKey(serviceName, req.UserId))

	info, err := GetUserInfo(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	//用户为冻结态
	if info.Status == 1 {
		return nil, UserFreezeError
	}

	//解绑邮箱
	if err := DB.Model(&models.User{}).Where("user_id", req.UserId).
		Update("updated_at", time.Now()).
		Update("email", "").Error; err != nil {
		util.LogError("解绑邮箱错误", "UnbindEmail", err)
		return nil, err
	}

	return &user.Empty{}, nil
}

// FreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) FreezeUser(ctx context.Context, req *user.FreezeUserReq) (resp *user.Empty, err error) {
	//清理缓存
	defer Rds.Del(ctx, util.TakeKey(serviceName, req.UserId))

	info, err := GetUserInfo(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	//状态已经修改
	if info.Status == 1 {
		return &user.Empty{}, nil
	}

	if err := DB.Model(&models.User{}).
		Where("id = ?", req.UserId).
		Update("status", 1).
		Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("冻结用户失败", "FreezeUser", err)
		return nil, err
	}

	return &user.Empty{}, nil
}

// UnfreezeUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UnfreezeUser(ctx context.Context, req *user.UnfreezeUserReq) (resp *user.Empty, err error) {
	//清理缓存
	defer Rds.Del(ctx, util.TakeKey(serviceName, req.UserId))

	info, err := GetUserInfo(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	//状态已经修改
	if info.Status == 0 {
		return &user.Empty{}, nil
	}

	if err := DB.Model(&models.User{}).
		Where("id = ?", req.UserId).
		Update("status", 0).
		Update("updated_at", time.Now()).Error; err != nil {
		util.LogError("解冻用户失败", "UnfreezeUser", err)
		return nil, err
	}

	return &user.Empty{}, nil
}
