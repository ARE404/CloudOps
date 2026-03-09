package startup

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/ai/service"
	"go.uber.org/zap"
)

// ApplicationBootstrap runs startup tasks
type ApplicationBootstrap struct {
	ragSvc service.RAGService
	logger *zap.Logger
}

func NewApplicationBootstrap(ragSvc service.RAGService, logger *zap.Logger) *ApplicationBootstrap {
	return &ApplicationBootstrap{ragSvc: ragSvc, logger: logger}
}

// Start loads the vector store and performs other startup tasks
func (b *ApplicationBootstrap) Start() {
	ctx := context.Background()
	if err := b.ragSvc.LoadVectorStore(ctx); err != nil {
		b.logger.Warn("向量知识库加载失败，AI RAG 功能可能不可用", zap.Error(err))
	}
}
