package main

import (
	"context"
	"errors"
	auth "github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var secretKey = config.Conf.Jwt.AdminSecretKey

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

type UserClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userId int) (string, error) {

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
		return "", err
	}

	return signedToken, nil
}

func ParseJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	token, err := GenerateJWT(int(req.UserId))
	if err != nil {
		klog.Fatal(err)
	}
	resp = &auth.DeliveryResp{Token: token}
	return resp, nil
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	token, err := ParseJWT(req.Token)
	resp = &auth.VerifyResp{Res: err == nil}
	if resp.Res {
		//如果相差时间小于令牌存活阈值,就重新生成令牌
		diff := config.Conf.AdminTtl - config.Conf.AdminSuv
		if diff < 0 {
			diff = 10800
		}
		suv := time.Duration(diff) * time.Second
		if time.Since(token.IssuedAt.Time) >= suv {
			newToken, err := GenerateJWT(token.UserId)
			if err != nil {
				return nil, err // 返回错误而不是直接终止程序
			}
			//将token重新放入
			resp.Token = newToken
		} else {
			resp.Token = req.Token
		}
	}
	return resp, err
}
