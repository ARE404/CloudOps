package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/GoSimplicity/CloudOps/internal/ai/dao"
	"github.com/GoSimplicity/CloudOps/internal/ai/vector"
	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
)

// RAGService handles knowledge base management and RAG queries
type RAGService interface {
	Chat(ctx context.Context, req *model.ChatReq) (*model.ChatResp, error)
	AddKnowledge(ctx context.Context, req *model.AddKnowledgeReq) error
	ListKnowledge(ctx context.Context, req *model.ListKnowledgeReq) (*model.PageResp, error)
	DeleteKnowledge(ctx context.Context, id int) error
	LoadVectorStore(ctx context.Context) error
}

type ragService struct {
	llm    LLMService
	dao    dao.KnowledgeDAO
	store  *vector.Store
	logger *zap.Logger
}

func NewRAGService(llm LLMService, dao dao.KnowledgeDAO, store *vector.Store, logger *zap.Logger) RAGService {
	return &ragService{llm: llm, dao: dao, store: store, logger: logger}
}

func (s *ragService) Chat(ctx context.Context, req *model.ChatReq) (*model.ChatResp, error) {
	var systemPrompt string

	if req.UseRAG {
		// Embed the query and retrieve relevant docs
		queryEmb, err := s.llm.Embed(ctx, req.Message)
		if err != nil {
			s.logger.Warn("向量检索失败，降级为直接对话", zap.Error(err))
		} else {
			chunks := s.store.Search(ctx, queryEmb, 3)
			if len(chunks) > 0 {
				var sb strings.Builder
				sb.WriteString("根据以下知识库内容回答用户问题：\n\n")
				for i, c := range chunks {
					sb.WriteString(fmt.Sprintf("【参考%d】%s\n%s\n\n", i+1, c.Title, c.Content))
				}
				systemPrompt = sb.String()
			}
		}
	}

	messages := make([]Message, 0, 3)
	if systemPrompt != "" {
		messages = append(messages, Message{Role: "system", Content: systemPrompt})
	} else {
		messages = append(messages, Message{Role: "system", Content: "你是一个专业的云运维助手，擅长 Kubernetes、Prometheus 和云原生技术。"})
	}
	messages = append(messages, Message{Role: "user", Content: req.Message})

	answer, err := s.llm.Chat(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &model.ChatResp{
		Answer:    answer,
		SessionID: req.SessionID,
	}, nil
}

func (s *ragService) AddKnowledge(ctx context.Context, req *model.AddKnowledgeReq) error {
	// Generate embedding
	emb, err := s.llm.Embed(ctx, req.Title+" "+req.Content)
	if err != nil {
		s.logger.Warn("生成 embedding 失败，仍保存文档", zap.Error(err))
	}

	embJSON := "[]"
	if len(emb) > 0 {
		data, _ := json.Marshal(emb)
		embJSON = string(data)
	}

	chunk := &model.AiKnowledgeChunk{
		Title:     req.Title,
		Content:   req.Content,
		Source:    req.Source,
		Category:  req.Category,
		Embedding: embJSON,
	}

	if err := s.dao.Create(ctx, chunk); err != nil {
		return err
	}

	// Update in-memory store
	if len(emb) > 0 {
		s.store.Add(&vector.Document{
			ID:        chunk.ID,
			Title:     chunk.Title,
			Content:   chunk.Content,
			Embedding: emb,
		})
	}
	return nil
}

func (s *ragService) ListKnowledge(ctx context.Context, req *model.ListKnowledgeReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *ragService) DeleteKnowledge(ctx context.Context, id int) error {
	if err := s.dao.Delete(ctx, id); err != nil {
		return err
	}
	s.store.Remove(id)
	return nil
}

func (s *ragService) LoadVectorStore(ctx context.Context) error {
	chunks, err := s.dao.ListAll(ctx)
	if err != nil {
		return err
	}
	s.store.LoadFromChunks(chunks)
	s.logger.Info("向量知识库加载完成", zap.Int("count", len(chunks)))
	return nil
}
