package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/db"
	"github.com/123508/douyinshop/pkg/models"
	"github.com/123508/douyinshop/pkg/myredis"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	serviceName = "user"
)

type UserClaims struct {
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

var secretKey = config.Conf.Jwt.AdminSecretKey

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", "", err)
	}
	return DB
}

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", "", err)
	}
	return rds
}

// GenerateJWT 产生一个jwt令牌
func GenerateJWT(userId uint64) (string, error) {

	claims := UserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.AdminTtl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("token签名失败:", err)
		return "", SignFailError
	}

	return signedToken, nil
}

// ParseJWT 解析jwt令牌
func ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("解析token失败:", err)
		return nil, ParseTokenError
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Println("token已过期")
		return nil, InvalidToken
	}
}

// Encryption sha256加密算法
func Encryption(origin string) string {
	hash := sha256.New()
	hash.Write([]byte(origin))
	hashBytes := hash.Sum(nil)
	res := hex.EncodeToString(hashBytes)
	return res
}

func GetUserInfo(ctx context.Context, userId uint64) (models.User, error) {
	var row models.User
	if err := DB.Model(&models.User{}).Where("id = ?", userId).First(&row).Error; err != nil {
		util.LogError("查询用户异常", "GetUserInfo", "", err)
		return models.User{}, SearchMySQLError
	}

	return row, nil

}

func GetUserPassword(ctx context.Context, userId uint64) (models.UserLogin, error) {
	var userLogin models.UserLogin
	if err := DB.Model(&models.UserLogin{}).Where("user_id = ?", userId).First(&userLogin).Error; err != nil {
		util.LogError("查询用户密码异常", "GetUserPassword", "", err)
		return models.UserLogin{}, err
	}
	return userLogin, nil
}
