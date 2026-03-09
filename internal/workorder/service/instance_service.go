package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/dao"
	"go.uber.org/zap"
)

type InstanceService interface {
	ListInstances(ctx context.Context, req *model.ListInstancesReq) (*model.PageResp, error)
	GetInstance(ctx context.Context, id int) (*model.WoInstance, error)
	CreateInstance(ctx context.Context, req *model.CreateInstanceReq) error
}

type instanceService struct {
	dao    dao.InstanceDAO
	logger *zap.Logger
}

func NewInstanceService(dao dao.InstanceDAO, logger *zap.Logger) InstanceService {
	return &instanceService{dao: dao, logger: logger}
}

func (s *instanceService) ListInstances(ctx context.Context, req *model.ListInstancesReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *instanceService) GetInstance(ctx context.Context, id int) (*model.WoInstance, error) {
	return s.dao.GetByID(ctx, id)
}

func (s *instanceService) CreateInstance(ctx context.Context, req *model.CreateInstanceReq) error {
	instance := &model.WoInstance{
		Title:        req.Title,
		TemplateID:   req.TemplateID,
		FormData:     req.FormData,
		Status:       "pending",
		Priority:     req.Priority,
		AssigneeID:   req.AssigneeID,
		ReporterID:   req.ReporterID,
		ReporterName: req.ReporterName,
		DueAt:        req.DueAt,
	}
	if instance.Priority == 0 {
		instance.Priority = 2
	}
	return s.dao.Create(ctx, instance)
}
