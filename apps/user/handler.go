package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/redis"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/123508/douyinshop/kitex_gen/user"
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

var DB = open()

func open() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	return DB
}

// Register implements the UserServiceImpl interface.
// 用户注册接口
// 如果两个密码不同,则返回空
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	if req.Password != req.ConfirmPassword {
		klog.Error("密码和注册密码不一致!")
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
		klog.Error("创建用户对象错误!")
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
		klog.Error("用户名或密码错误")
		return &user.LoginResp{UserId: 0}, ErrorUsernameOrPassword
	}

	//这个方法之后要调用发放token的逻辑

	return &user.LoginResp{UserId: row.ID}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
// 获取用户信息接口
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {

	var row models.User
	DB.Model(&models.User{}).Where("id = ?", req.UserId).First(&row)

	//如果查询到不存在该用户,返回error
	if row.Email == "" {
		klog.Error("用户不存在")
		return nil, UserNotExists
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
func (s *UserServiceImpl) Logout(ctx context.Context, req *user.LogoutReq) (resp *user.LogoutResp, err error) {

	ir, err := redis.InitRedis()

	//如果初始化redis成功则让token失效,否则报错并返回
	if err != nil {
		klog.Error("初始化redis错误")
		//如果redis初始化错误则返回nil
		return nil, err
	} else {
		ir.Set(ctx, req.Token, "1", 8*time.Hour)
		return &user.LogoutResp{}, nil
	}
}

// Update implements the UserServiceImpl interface.
// 用户更新接口
// 更新用户资料并
// 绑定整个更新为事务,如果出错就进行回滚
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.UpdateResp, err error) {

	info, err := s.GetUserInfo(ctx, &user.GetUserInfoReq{UserId: req.UserId})

	if err != nil {
		klog.Error("用户不存在")
		return nil, UserNotExists
	}

	err = DB.Transaction(func(tx *gorm.DB) error {

		if req.UserId == 0 {
			klog.Error("用户不存在")
			return UserNotExists
		}

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
			klog.Error("更新用户发生错误")
			return err
		}

		//更新用户密码部分
		if req.Password != "" {
			where := DB.Model(&models.UserLogin{}).Where("user_id = ?", req.UserId)
			if err = where.Update("password", encryption(req.Password)).Update("updated_at", time.Now()).Error; err != nil {
				klog.Error("更新用户密码错误")
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
