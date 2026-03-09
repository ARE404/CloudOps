package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/dao"
	"go.uber.org/zap"
)

type CommentService interface {
	ListComments(ctx context.Context, instanceID int) ([]*model.WoComment, error)
	CreateComment(ctx context.Context, req *model.CreateCommentReq) error
}

type commentService struct {
	dao    dao.CommentDAO
	logger *zap.Logger
}

func NewCommentService(dao dao.CommentDAO, logger *zap.Logger) CommentService {
	return &commentService{dao: dao, logger: logger}
}

func (s *commentService) ListComments(ctx context.Context, instanceID int) ([]*model.WoComment, error) {
	return s.dao.ListByInstance(ctx, instanceID)
}

func (s *commentService) CreateComment(ctx context.Context, req *model.CreateCommentReq) error {
	comment := &model.WoComment{
		InstanceID: req.InstanceID,
		UserID:     req.UserID,
		Username:   req.Username,
		Content:    req.Content,
	}
	return s.dao.Create(ctx, comment)
}
