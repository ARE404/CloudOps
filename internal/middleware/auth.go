package middleware

import (
	"strings"

	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var noAuthPaths = []string{
	"/api/auth/login",
	"/api/auth/register",
	"/swagger/",
	"/health",
}

type JWTMiddleware struct {
	jwtHandler jwt.Handler
}

func NewJWTMiddleware(h jwt.Handler) *JWTMiddleware {
	return &JWTMiddleware{jwtHandler: h}
}

func (m *JWTMiddleware) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, p := range noAuthPaths {
			if strings.HasPrefix(path, p) {
				ctx.Next()
				return
			}
		}

		tokenStr := m.jwtHandler.ExtractToken(ctx)
		if tokenStr == "" {
			base.UnauthorizedError(ctx, "请先登录")
			ctx.Abort()
			return
		}

		key := []byte(viper.GetString("jwt.key1"))
		claims := &jwt.UserClaims{}
		token, err := gojwt.ParseWithClaims(tokenStr, claims, func(t *gojwt.Token) (interface{}, error) {
			return key, nil
		})
		if err != nil || !token.Valid {
			base.UnauthorizedError(ctx, "token 无效或已过期")
			ctx.Abort()
			return
		}

		if err := m.jwtHandler.CheckSession(ctx, claims.Ssid); err != nil {
			base.UnauthorizedError(ctx, "会话已失效，请重新登录")
			ctx.Abort()
			return
		}

		ctx.Set("user", *claims)
		ctx.Next()
	}
}
