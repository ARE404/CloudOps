package base

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	StatusSuccess = 0
	StatusError   = 1
)

func apiData(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(http.StatusOK, ApiResponse{Code: code, Data: data, Message: message})
}

func Success(c *gin.Context) {
	apiData(c, StatusSuccess, map[string]interface{}{}, "操作成功")
}

func SuccessWithData(c *gin.Context, data interface{}) {
	apiData(c, StatusSuccess, data, "请求成功")
}

func ErrorWithMessage(c *gin.Context, message string) {
	apiData(c, StatusError, map[string]interface{}{}, message)
}

func BadRequestError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ApiResponse{Code: StatusError, Data: map[string]interface{}{}, Message: message})
}

func BadRequestWithDetails(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusBadRequest, ApiResponse{Code: StatusError, Data: data, Message: message})
}

func UnauthorizedError(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ApiResponse{Code: StatusError, Data: map[string]interface{}{}, Message: message})
}

func ForbiddenError(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ApiResponse{Code: StatusError, Data: map[string]interface{}{}, Message: message})
}

// HandleRequest binds request, executes action, and writes response.
func HandleRequest(ctx *gin.Context, req interface{}, action func() (interface{}, error)) {
	if req != nil {
		if err := ctx.ShouldBind(req); err != nil {
			BadRequestWithDetails(ctx, err.Error(), "参数绑定失败")
			return
		}
	}
	data, err := action()
	if err != nil {
		ErrorWithMessage(ctx, err.Error())
		return
	}
	if data != nil {
		SuccessWithData(ctx, data)
	} else {
		Success(ctx)
	}
}

func GetCustomParamID(ctx *gin.Context, paramName string) (int, error) {
	val := ctx.Param(paramName)
	if val == "" {
		return 0, fmt.Errorf("缺少路径参数 '%s'", paramName)
	}
	id, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("路径参数 '%s' 不是整数", paramName)
	}
	return id, nil
}

func GetParamID(ctx *gin.Context) (int, error) {
	return GetCustomParamID(ctx, "id")
}
