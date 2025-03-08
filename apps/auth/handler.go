package main

import (
	"context"
	"errors"
	"github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/123508/douyinshop/pkg/errorno"
	"github.com/123508/douyinshop/pkg/redis"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

var secretKey = config.Conf.Jwt.AdminSecretKey

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

type UserClaims struct {
	UserId uint32 `json:"user_id"`
	jwt.RegisteredClaims
}

var InvalidToken = &errorno.BasicMessageError{Code: 401, Message: "token已过期,请重新登录"}

var ParseTokenError = &errorno.BasicMessageError{Code: 402, Message: "解析token失败"}

var SignFailError = &errorno.BasicMessageError{Code: 500, Message: "token签名失败"}

var RedisError = &errorno.BasicMessageError{Code: 500, Message: "Redis数据库异常"}

var RedisConnectionError = &errorno.BasicMessageError{Code: 404, Message: "Redis数据库连接异常"}

func GenerateJWT(userId uint32) (string, error) {

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

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	token, err := GenerateJWT(req.UserId)
	if err != nil {
		return nil, err
	}
	resp = &auth.DeliveryResp{Token: token}
	return resp, nil
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
// 验证令牌接口
// 如果redis中标记该令牌无效,返回错误响应
// 如果令牌无效,返回错误响应
// 如果令牌存活时间小于等于阈值,刷新令牌并返回成功响应
// 如果令牌存活时间大于阈值,直接返回成功响应
// 注意每次需要使用响应去接收token
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {

	if req.Token == "" {
		return &auth.VerifyResp{Res: false}, errors.New("请求令牌为空")
	}

	rds, err := redis.InitRedis()
	if err != nil {
		klog.Fatal(err)
		return nil, RedisError
	}
	//在redis中检查token是否存活
	result, err := rds.Exists(ctx, req.Token).Result()

	//如果redis连接出错,直接返回错误信息
	if err != nil {
		log.Println("redis连接错误:", err)
		return nil, RedisConnectionError
	} else {
		//如果在redis中检测到token,则直接返回失败响应
		if result == 1 {
			resp = &auth.VerifyResp{Res: false}
			return resp, InvalidToken
		}
	}

	token, err := ParseJWT(req.Token)
	//判断令牌是否可以被解析,如果令牌无法被解析返回失败响应
	resp = &auth.VerifyResp{Res: err == nil}

	if resp.Res {
		//如果相差时间小于令牌存活阈值,就重新生成令牌
		diff := config.Conf.AdminTtl - config.Conf.AdminSuv
		if diff <= 0 {
			diff = 10800
		}
		suv := time.Duration(diff) * time.Second
		if time.Since(token.IssuedAt.Time) >= suv {
			newToken, err := GenerateJWT(token.UserId)
			if err != nil {
				resp.Res = false
				return resp, err // 返回错误
			}
			//将token重新放入
			resp.Token = newToken
		} else {
			resp.Token = req.Token
		}
		resp.UserId = token.UserId
	}

	return resp, err
}
