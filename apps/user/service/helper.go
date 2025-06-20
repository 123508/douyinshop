package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/util"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

var backendSecretKey = config.Conf.Jwt.ServiceSecretKey

// GenerateFrontendJWT 产生一个jwt令牌
func GenerateFrontendJWT(userId uint64) (string, error) {

	jti, _ := uuid.NewV7()

	claims := FrontendClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Conf.AdminTtl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    util.Md5Hash("user service delivery this token"),
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	fmt.Println(backendSecretKey)

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
