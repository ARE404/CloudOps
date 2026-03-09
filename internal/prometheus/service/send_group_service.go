package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/dao"
	"go.uber.org/zap"
)

type SendGroupService interface {
	ListSendGroups(ctx context.Context, req *model.ListSendGroupsReq) (*model.PageResp, error)
	CreateSendGroup(ctx context.Context, req *model.CreateSendGroupReq) error
	UpdateSendGroup(ctx context.Context, req *model.UpdateSendGroupReq) error
	DeleteSendGroup(ctx context.Context, id int) error
}

type sendGroupService struct {
	dao    dao.SendGroupDAO
	logger *zap.Logger
}

func NewSendGroupService(dao dao.SendGroupDAO, logger *zap.Logger) SendGroupService {
	return &sendGroupService{dao: dao, logger: logger}
}

func (s *sendGroupService) ListSendGroups(ctx context.Context, req *model.ListSendGroupsReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *sendGroupService) CreateSendGroup(ctx context.Context, req *model.CreateSendGroupReq) error {
	group := &model.PromSendGroup{
		Name:           req.Name,
		PoolID:         req.PoolID,
		Channels:       req.Channels,
		RepeatInterval: req.RepeatInterval,
		Description:    req.Description,
	}
	if group.RepeatInterval == 0 {
		group.RepeatInterval = 3600
	}
	return s.dao.Create(ctx, group)
}

func (s *sendGroupService) UpdateSendGroup(ctx context.Context, req *model.UpdateSendGroupReq) error {
	group, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Channels != "" {
		group.Channels = req.Channels
	}
	if req.RepeatInterval > 0 {
		group.RepeatInterval = req.RepeatInterval
	}
	if req.Description != "" {
		group.Description = req.Description
	}
	return s.dao.Update(ctx, group)
}

func (s *sendGroupService) DeleteSendGroup(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
