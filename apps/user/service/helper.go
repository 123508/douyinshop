package service

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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

const (
	serviceName = "user"
)

type FrontendClaims struct {
	UserId    uint64 `json:"user_id"`
	SessionId uint64 `json:"session_id"`
	jwt.RegisteredClaims
}

type BackendClaims struct {
	UserId uint64   `json:"user_id"`
	Roles  []string `json:"roles"`
	Perms  []string `json:"perms"`
	PVer   uint64   `json:"p_ver"`
	jwt.RegisteredClaims
}

var frontendSecretKey = config.Conf.Jwt.AdminSecretKey

var backendSecretKey = config.Conf.Jwt.AdminSecretKey

var DB = connectWithMySQL()

func connectWithMySQL() *gorm.DB {
	DB, err := db.InitDB()
	if err != nil {
		util.LogError("打开MySQL连接失败", "connectWithMySQL", err)
	}
	return DB
}

var Rds = connectWithRedis()

func connectWithRedis() *redis.Client {
	rds, err := myredis.InitRedis()
	if err != nil {
		util.LogError("打开Redis连接失败", "connectWithRedis", err)
	}
	return rds
}

// GenerateFrontendJWT 产生一个jwt令牌
func GenerateFrontendJWT(userId uint64) (string, error) {

	jti, _ := uuid.NewV7()

	claims := FrontendClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.AdminTtl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "user service delivery this token",
			ID:        jti.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(frontendSecretKey))
	if err != nil {
		util.LogError("token签名失败:", "GenerateFrontendJWT", err)
		return "", SignFailError
	}

	return signedToken, nil
}

// ParseFrontendJWT 解析jwt令牌
func ParseFrontendJWT(tokenString string) (*FrontendClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &FrontendClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(frontendSecretKey), nil
	})

	if err != nil {
		util.LogError("解析token失败", "ParseFrontendJWT", err)
		return nil, ParseTokenError
	}

	if claims, ok := token.Claims.(*FrontendClaims); ok && token.Valid {
		return claims, nil
	} else {
		util.LogError("token已过期", "ParseFrontendJWT", err)
		return nil, TokenTimeOutError
	}
}

func GenerateBackendJWT(
	userId uint64,
	roles []string,
	perms []string,
	version uint64,
) (string, error) {

	jti, _ := uuid.NewV7()

	claims := BackendClaims{
		UserId: userId,
		Roles:  roles,
		Perms:  perms,
		PVer:   version,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.AdminTtl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    util.Md5Hash("user service delivery this token"),
			ID:        jti.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	signedToken, err := token.SignedString([]byte(backendSecretKey))
	if err != nil {
		util.LogError("token签名失败:", "GenerateBackendJWT", err)
		return "", SignFailError
	}

	return signedToken, nil
}

func ParseBackendJWT(tokenString string) (*BackendClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &BackendClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(backendSecretKey), nil
	})
	if err != nil {
		return nil, ParseTokenError
	}

	if claims, ok := token.Claims.(*BackendClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, TokenTimeOutError
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
		util.LogError("查询用户异常", "GetUserInfo", err)
		return models.User{}, SearchMySQLError
	}

	return row, nil

}
