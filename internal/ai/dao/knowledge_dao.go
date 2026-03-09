package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KnowledgeDAO interface {
	List(ctx context.Context, req *model.ListKnowledgeReq) ([]*model.AiKnowledgeChunk, int64, error)
	Create(ctx context.Context, chunk *model.AiKnowledgeChunk) error
	Delete(ctx context.Context, id int) error
	ListAll(ctx context.Context) ([]*model.AiKnowledgeChunk, error)
}

type knowledgeDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewKnowledgeDAO(db *gorm.DB, logger *zap.Logger) KnowledgeDAO {
	return &knowledgeDAO{db: db, logger: logger}
}

func (d *knowledgeDAO) List(ctx context.Context, req *model.ListKnowledgeReq) ([]*model.AiKnowledgeChunk, int64, error) {
	var chunks []*model.AiKnowledgeChunk
	var total int64
	query := d.db.WithContext(ctx).Model(&model.AiKnowledgeChunk{}).Where("deleted_at IS NULL")
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Select("id, title, source, category, created_at, updated_at").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&chunks).Error; err != nil {
		return nil, 0, err
	}
	return chunks, total, nil
}

func (d *knowledgeDAO) Create(ctx context.Context, chunk *model.AiKnowledgeChunk) error {
	return d.db.WithContext(ctx).Create(chunk).Error
}

func (d *knowledgeDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.AiKnowledgeChunk{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (d *knowledgeDAO) ListAll(ctx context.Context) ([]*model.AiKnowledgeChunk, error) {
	var chunks []*model.AiKnowledgeChunk
	if err := d.db.WithContext(ctx).Where("deleted_at IS NULL AND embedding != ''").Find(&chunks).Error; err != nil {
		return nil, err
	}
	return chunks, nil
}
