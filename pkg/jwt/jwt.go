package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const (
	bearerPrefix             = "Bearer "
	authorizationHeader      = "Authorization"
	sessionKeyPattern        = "cloudops:user:ssid:%s"
	tokenBlacklistKeyPattern = "cloudops:blacklist:token:%s"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int, username string, accountType int8) (string, string, error)
	SetJWTToken(ctx *gin.Context, uid int, username string, ssid string, accountType int8) (string, error)
	ExtractToken(ctx *gin.Context) string
	CheckSession(ctx *gin.Context, ssid string) error
	ClearToken(ctx *gin.Context) error
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid         int
	Username    string
	Ssid        string
	AccountType int8
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Uid         int
	Username    string
	Ssid        string
	AccountType int8
}

type handler struct {
	client        redis.Cmdable
	signingMethod jwt.SigningMethod
	jwtExpiration time.Duration
	rcExpiration  time.Duration
	key1          []byte
	key2          []byte
	issuer        string
}

func NewJWTHandler(c redis.Cmdable) Handler {
	key1 := viper.GetString("jwt.key1")
	key2 := viper.GetString("jwt.key2")
	issuer := viper.GetString("jwt.issuer")
	expirationMinutes := viper.GetInt64("jwt.expiration")
	if expirationMinutes <= 0 {
		expirationMinutes = 1440
	}
	return &handler{
		client:        c,
		signingMethod: jwt.SigningMethodHS512,
		jwtExpiration: time.Minute * time.Duration(expirationMinutes),
		rcExpiration:  time.Hour * 24 * 7,
		key1:          []byte(key1),
		key2:          []byte(key2),
		issuer:        issuer,
	}
}

func (h *handler) SetLoginToken(ctx *gin.Context, uid int, username string, accountType int8) (string, string, error) {
	ssid := uuid.New().String()
	refreshToken, err := h.setRefreshToken(uid, username, ssid, accountType)
	if err != nil {
		return "", "", err
	}
	jwtToken, err := h.SetJWTToken(ctx, uid, username, ssid, accountType)
	if err != nil {
		return "", "", err
	}
	return jwtToken, refreshToken, nil
}

func (h *handler) SetJWTToken(_ *gin.Context, uid int, username string, ssid string, accountType int8) (string, error) {
	uc := UserClaims{
		Uid:         uid,
		Username:    username,
		Ssid:        ssid,
		AccountType: accountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtExpiration)),
			Issuer:    h.issuer,
		},
	}
	return h.sign(uc, h.key1)
}

func (h *handler) setRefreshToken(uid int, username string, ssid string, accountType int8) (string, error) {
	rc := RefreshClaims{
		Uid:         uid,
		Username:    username,
		Ssid:        ssid,
		AccountType: accountType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.rcExpiration)),
		},
	}
	return h.sign(rc, h.key2)
}

func (h *handler) ExtractToken(ctx *gin.Context) string {
	authHeader := strings.TrimSpace(ctx.GetHeader(authorizationHeader))
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(authHeader[len(bearerPrefix):])
}

func (h *handler) CheckSession(ctx *gin.Context, ssid string) error {
	c, err := h.client.Exists(ctx, fmt.Sprintf(sessionKeyPattern, ssid)).Result()
	if err != nil {
		return err
	}
	if c != 0 {
		return errors.New("token 已失效")
	}
	return nil
}

func (h *handler) ClearToken(ctx *gin.Context) error {
	token := h.ExtractToken(ctx)
	if token == "" {
		return errors.New("missing token")
	}
	claims := &UserClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return h.key1, nil
	})
	if err != nil || !t.Valid {
		return errors.New("invalid token")
	}
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		remaining = time.Second
	}
	return h.client.Set(ctx, fmt.Sprintf(tokenBlacklistKeyPattern, token), "invalid", remaining).Err()
}

func (h *handler) sign(claims jwt.Claims, key []byte) (string, error) {
	t := jwt.NewWithClaims(h.signingMethod, claims)
	return t.SignedString(key)
}
