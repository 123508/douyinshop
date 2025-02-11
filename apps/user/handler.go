package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/redis"
	"github.com/123508/douyinshop/pkg/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"log"
	"time"

	user "github.com/123508/douyinshop/kitex_gen/user"
)

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
// 如果两个密码不同,则返回0用户
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	if req.Password != req.ConfirmPassword {
		klog.Fatal("密码不匹配,请重新输入")
		return &user.RegisterResp{UserId: 0}, nil
	}
	//初始化用户属性
	user1 := &models.User{}
	user1.Name = req.Nickname
	user1.Email = req.Email
	user1.Phone = req.Phone
	user1.Gender = req.Gender
	user1.CreatedAt = time.Now()
	user1.UpdatedAt = time.Now()
	user1.Avatar = utils.UploadImages("", "user", 0)
	user1.IsLive = 0
	//创建用户对象
	DB.Create(user1)
	user2 := &models.UserLogin{}
	user2.UserId = uint32(user1.ID)
	user2.User = *user1
	user2.Password = encryption(req.Password)
	user2.CreatedAt = time.Now()
	user2.UpdatedAt = time.Now()
	user2.IsLive = 0
	DB.Create(user2)
	return &user.RegisterResp{UserId: uint32(user1.ID)}, nil
}

// Login implements the UserServiceImpl interface.
// 用户登录接口
// 如果账号密码正确，返回用户ID
// 如果账号或密码错误，user_id为0
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	var row models.User
	DB.Model(&models.User{}).Where("email = ?", req.Email).Find(&row)
	var res models.UserLogin
	DB.Model(&models.UserLogin{}).Where("user_id = ? and password = ?", row.ID, encryption(req.Password)).Find(&res)
	if res.ID == 0 {
		return &user.LoginResp{UserId: 0}, nil
	}

	return &user.LoginResp{UserId: uint32(row.ID)}, nil
}

// GetUserInfo implements the UserServiceImpl interface.
// 获取用户信息接口
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {

	var row models.User
	DB.Model(&models.User{}).Where("id = ?", req.UserId).Find(&row)
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
func (s *UserServiceImpl) Logout(ctx context.Context, req *user.LogoutReq) (resp *user.LogoutResp, err error) {

	ir, err := redis.InitRedis()
	if err != nil {
		klog.Fatal(err)
	}
	ir.Set(ctx, req.Token, "1", 8*time.Hour)

	return &user.LogoutResp{}, nil
}

// Update implements the UserServiceImpl interface.
// 用户更新接口
// 更新用户资料并
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.UpdateResp, err error) {
	images := utils.UploadImages(req.Avatar, "user", req.UserId)
	tx := DB.Model(&models.User{}).Where("id=?", req.UserId).Update("gender", req.Gender)
	tx.Update("phone", req.Phone).Update("avatar", images).Update("name", req.Nickname)
	tx.Update("update_time", time.Now())

	if req.Password != "" {
		where := DB.Model(&models.UserLogin{}).Where("user_id = ?", req.UserId)
		where.Update("password", encryption(req.Password)).Update("update_time", time.Now())
	}
	return &user.UpdateResp{}, nil
}

// Delete implements the UserServiceImpl interface.
// 删除用户接口
func (s *UserServiceImpl) Delete(ctx context.Context, req *user.DeleteReq) (resp *user.DeleteResp, err error) {
	DB.Model(&models.User{}).Where("id = ?", req.UserId).Update("is_live", 1)
	DB.Model(&models.UserLogin{}).Where("user_id=?", req.UserId).Update("is_live", 1)
	return &user.DeleteResp{}, nil
}
