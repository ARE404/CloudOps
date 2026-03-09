package api

import (
	"net/http"

	"github.com/GoSimplicity/CloudOps/internal/ai/service"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/pkg/base"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	ragSvc service.RAGService
}

func NewAIHandler(ragSvc service.RAGService) *AIHandler {
	return &AIHandler{ragSvc: ragSvc}
}

func (h *AIHandler) RegisterRouters(server *gin.Engine) {
	g := server.Group("/api/ai")
	{
		g.GET("/health", h.Health)
		g.POST("/chat", h.Chat)
		g.GET("/knowledge", h.ListKnowledge)
		g.POST("/knowledge", h.AddKnowledge)
		g.DELETE("/knowledge/:id", h.DeleteKnowledge)
	}
}

func (h *AIHandler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ai"})
}

func (h *AIHandler) Chat(ctx *gin.Context) {
	var req model.ChatReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.ragSvc.Chat(ctx, &req)
	})
}

func (h *AIHandler) ListKnowledge(ctx *gin.Context) {
	var req model.ListKnowledgeReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return h.ragSvc.ListKnowledge(ctx, &req)
	})
}

func (h *AIHandler) AddKnowledge(ctx *gin.Context) {
	var req model.AddKnowledgeReq
	base.HandleRequest(ctx, &req, func() (interface{}, error) {
		return nil, h.ragSvc.AddKnowledge(ctx, &req)
	})
}

func (h *AIHandler) DeleteKnowledge(ctx *gin.Context) {
	id, err := base.GetParamID(ctx)
	if err != nil {
		base.BadRequestError(ctx, err.Error())
		return
	}
	base.HandleRequest(ctx, nil, func() (interface{}, error) {
		return nil, h.ragSvc.DeleteKnowledge(ctx, id)
	})
}
