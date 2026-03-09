package api

import (
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/service"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	svc service.CommentService
}

func NewCommentHandler(svc service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

func (h *CommentHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/workorder/instances/:id/comments")
	{
		g.GET("", h.List)
		g.POST("", h.Create)
	}
}

func (h *CommentHandler) List(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return h.svc.ListComments(ctx, id)
	})
}

func (h *CommentHandler) Create(ctx *gin.Context) {
	var req model.CreateCommentReq
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	req.InstanceID = id
	uc := ctx.MustGet("user").(jwt.UserClaims)
	req.UserID = uc.Uid
	req.Username = uc.Username
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.svc.CreateComment(ctx, &req)
	})
}
