package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentDAO interface {
	ListByInstance(ctx context.Context, instanceID int) ([]*model.WoComment, error)
	Create(ctx context.Context, comment *model.WoComment) error
}

type commentDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewCommentDAO(db *gorm.DB, logger *zap.Logger) CommentDAO {
	return &commentDAO{db: db, logger: logger}
}

func (d *commentDAO) ListByInstance(ctx context.Context, instanceID int) ([]*model.WoComment, error) {
	var comments []*model.WoComment
	if err := d.db.WithContext(ctx).Where("instance_id = ?", instanceID).Order("created_at ASC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (d *commentDAO) Create(ctx context.Context, comment *model.WoComment) error {
	return d.db.WithContext(ctx).Create(comment).Error
}
