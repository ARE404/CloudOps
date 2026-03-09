package di

import (
	"github.com/GoSimplicity/CloudOps/internal/middleware"
	ijwt "github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitMiddlewares(jwtHandler ijwt.Handler, logger *zap.Logger, enforcer *casbin.Enforcer) []gin.HandlerFunc {
	handlers := []gin.HandlerFunc{
		middleware.NewCORSMiddleware(),
		middleware.NewJWTMiddleware(jwtHandler).CheckLogin(),
		middleware.NewLogMiddleware(logger).Log(),
	}
	if enforcer != nil {
		handlers = append(handlers, middleware.NewCasbinMiddleware(enforcer).CheckAuth())
	}
	return handlers
}
