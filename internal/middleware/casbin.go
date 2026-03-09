package middleware

import (
	"fmt"

	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

func NewCasbinMiddleware(enforcer *casbin.Enforcer) *CasbinMiddleware {
	return &CasbinMiddleware{enforcer: enforcer}
}

func (m *CasbinMiddleware) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userVal, exists := ctx.Get("user")
		if !exists {
			ctx.Next()
			return
		}

		claims, ok := userVal.(jwt.UserClaims)
		if !ok {
			ctx.Next()
			return
		}

		// admin bypass
		if claims.AccountType == 2 {
			ctx.Next()
			return
		}

		sub := fmt.Sprintf("user:%d", claims.Uid)
		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		ok, err := m.enforcer.Enforce(sub, obj, act)
		if err != nil || !ok {
			base.ForbiddenError(ctx, "无权限访问该资源")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
