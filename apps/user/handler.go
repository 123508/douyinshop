package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/123508/douyinshop/kitex_gen/user"
)

const (
	serviceName = "user"
)

var UserNotExists = &errorno.BasicMessageError{Code: 401, Message: "用户不存在"}

var PasswordNotEqual = &errorno.BasicMessageError{Code: 400, Message: "密码不匹配,请重新输入"}

var ErrorUsernameOrPassword = &errorno.BasicMessageError{Code: 404, Message: "用户名或密码错误"}

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// sha256加密算法
func encryption(origin string) string {
	hash := sha256.New()
	hash.Write([]byte(origin))
	hashBytes := hash.Sum(nil)
	res := hex.EncodeToString(hashBytes)
	return res
}

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", "", err)
	}
	return DB
}

var rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

func GetUserInfoWithCache(ctx context.Context, userId uint64) (models.User, error) {
	simple := util.SimpleCacheComponent[uint64, models.User]{
		Rds:       rds,
		Ctx:       ctx,
		Key:       util.TakeKey(serviceName, userId),
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
		QueryExec: func() (models.User, error) {
			var row models.User
			if err := DB.Model(&models.User{}).Where("id = ?", userId).First(&row).Error; err != nil {
				util.LogError("用户不存在", "GetUserInfo", "", err)
				return models.User{}, UserNotExists
			}
			return row, nil
		},
		Expires: time.Duration(rand.Intn(10)+5) * time.Minute,
	}
	return simple.QueryWithCache()
}

// Register implements the UserServiceImpl interface.
// 用户注册接口
// 如果两个密码不同,则返回空
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	if req.Password != req.ConfirmPassword {
		util.LogError("密码和注册密码不一致", "Register", "密码比对", nil)
		return nil, PasswordNotEqual
	}

	//初始化用户属性
	user1 := &models.User{}
	user1.Email = req.Email
	user1.Name = req.Nickname
	user1.Phone = req.Phone
	user1.Gender = req.Gender
	user2 := &models.UserLogin{}
	user2.Password = encryption(req.Password)

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
		util.LogError("创建用户对象错误", "Register", "注册提交", err)
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
	DB.Model(&models.UserLogin{}).Where("user_id = ? and password = ?", row.ID, encryption(req.Password)).Find(&res)

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

	row, err := GetUserInfoWithCache(ctx, req.UserId)

	if err != nil {
		return nil, err
	}

	////通用类型
	//var row models.User
	//
	////查询缓存
	//result, _ := rds.Get(ctx, util.TakeKey(serviceName, req.UserId)).Result()
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
	//	err := rds.Set(ctx, util.TakeKey(serviceName, req.UserId), string(jsonData), time.Duration(rand.Intn(15)+30)*time.Minute).Err()
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
func (s *UserServiceImpl) Logout(ctx context.Context, req *user.LogoutReq) (resp *user.LogoutResp, err error) {

	//让token失效,否则报错并返回
	if err = rds.Set(ctx, req.Token, "1", 8*time.Hour).Err(); err != nil {
		return &user.LogoutResp{}, err
	}

	return &user.LogoutResp{}, nil

}

// Update implements the UserServiceImpl interface.
// 用户更新接口
// 更新用户资料并
// 绑定整个更新为事务,如果出错就进行回滚
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.UpdateResp, err error) {

	//清理缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, req.UserId))

	if req.UserId == 0 {
		return nil, UserNotExists
	}

	info, err := GetUserInfoWithCache(ctx, req.UserId)

	if err != nil {
		return nil, UserNotExists
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		//更新用户信息部分
		tx = DB.Model(&models.User{}).Where("id=?", req.UserId)

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
		if err := tx.Updates(updates).Error; err != nil {
			util.LogError("更新用户出错", "Update", "", err)
			return err
		}

		//更新用户密码部分
		if req.Password != "" {
			where := DB.Model(&models.UserLogin{}).Where("user_id = ?", req.UserId)
			if err = where.Update("password", encryption(req.Password)).Update("updated_at", time.Now()).Error; err != nil {
				util.LogError("更新用户密码错误", "Update", "", err)
				return err
			}
		}
		//返回nil提交事务
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user.UpdateResp{}, nil
}

// Delete implements the UserServiceImpl interface.
// 删除用户接口
// 绑定事务
func (s *UserServiceImpl) Delete(ctx context.Context, req *user.DeleteReq) (resp *user.DeleteResp, err error) {

	//清理缓存
	defer util.CleanCache(rds, ctx, util.TakeKey(serviceName, req.UserId))

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
	return &user.DeleteResp{}, nil
}
