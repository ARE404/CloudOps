package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/dao"
	"go.uber.org/zap"
)

type TemplateService interface {
	ListTemplates(ctx context.Context, req *model.ListTemplatesReq) (*model.PageResp, error)
	CreateTemplate(ctx context.Context, req *model.CreateTemplateReq) error
	UpdateTemplate(ctx context.Context, req *model.UpdateTemplateReq) error
	DeleteTemplate(ctx context.Context, id int) error
}

type templateService struct {
	dao    dao.TemplateDAO
	logger *zap.Logger
}

func NewTemplateService(dao dao.TemplateDAO, logger *zap.Logger) TemplateService {
	return &templateService{dao: dao, logger: logger}
}

func (s *templateService) ListTemplates(ctx context.Context, req *model.ListTemplatesReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *templateService) CreateTemplate(ctx context.Context, req *model.CreateTemplateReq) error {
	t := &model.WoTemplate{
		Name:        req.Name,
		Description: req.Description,
		FormSchema:  req.FormSchema,
		Category:    req.Category,
		Status:      1,
	}
	return s.dao.Create(ctx, t)
}

func (s *templateService) UpdateTemplate(ctx context.Context, req *model.UpdateTemplateReq) error {
	t, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if req.FormSchema != "" {
		t.FormSchema = req.FormSchema
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	return s.dao.Update(ctx, t)
}

func (s *templateService) DeleteTemplate(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
