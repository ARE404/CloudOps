package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Refresh-Token"},
		ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost") || strings.Contains(origin, "cloudops")
		},
		MaxAge: 12 * time.Hour,
	})
}
