package middleware

import (
	"context"
	"github.com/123508/douyinshop/apps/api/infras/client"
	auth "github.com/123508/douyinshop/kitex_gen/auth"
	"github.com/123508/douyinshop/pkg/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// ParseToken 解析token
func ParseToken() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := string(c.Cookie("token"))

		// 解析jwt
		userId, newToken, err := client.VerifyToken(ctx, &auth.VerifyTokenReq{Token: token})
		if err != nil || userId == 0 {
			c.JSON(401, map[string]interface{}{
				"error": "请先登录",
			})
			c.Abort()
			return
		}

		// 如果token发生变化，更新cookie
		if token != newToken {
			c.SetCookie(
				"token",
				newToken,
				config.Conf.Jwt.AdminTtl,
				"/",
				"localhost",
				protocol.CookieSameSiteLaxMode,
				false,
				true,
			)
		}
		c.Next(context.WithValue(ctx, "userId", userId))
	}
}
