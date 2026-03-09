package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) RegisterRouters(server *gin.Engine) {
	auth := server.Group("/api/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
	}
	users := server.Group("/api/users")
	{
		users.GET("", h.ListUsers)
		users.GET("/profile", h.GetProfile)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req model.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	resp, err := h.userSvc.Login(ctx, &req)
	if err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}
	base.SuccessWithData(ctx, resp)
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	if err := h.userSvc.Logout(ctx); err != nil {
		base.ErrorWithMessage(ctx, err.Error())
		return
	}
	base.Success(ctx)
}

func (h *UserHandler) ListUsers(ctx *gin.Context) {
	var req model.ListUsersReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.userSvc.ListUsers(ctx, &req)
	})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	uc := ctx.MustGet("user").(jwt.UserClaims)
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.userSvc.GetProfile(ctx, uc.Uid)
	})
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req model.CreateUserReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.userSvc.CreateUser(ctx, &req)
	})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req model.UpdateUserReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.ID = id
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.userSvc.UpdateUser(ctx, &req)
	})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.userSvc.DeleteUser(ctx, id)
	})
}
