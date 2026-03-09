package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/dao"
	"go.uber.org/zap"
)

type AlertRuleService interface {
	ListAlertRules(ctx context.Context, req *model.ListAlertRulesReq) (*model.PageResp, error)
	CreateAlertRule(ctx context.Context, req *model.CreateAlertRuleReq) error
	UpdateAlertRule(ctx context.Context, req *model.UpdateAlertRuleReq) error
	DeleteAlertRule(ctx context.Context, id int) error
}

type alertRuleService struct {
	dao    dao.AlertRuleDAO
	logger *zap.Logger
}

func NewAlertRuleService(dao dao.AlertRuleDAO, logger *zap.Logger) AlertRuleService {
	return &alertRuleService{dao: dao, logger: logger}
}

func (s *alertRuleService) ListAlertRules(ctx context.Context, req *model.ListAlertRulesReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *alertRuleService) CreateAlertRule(ctx context.Context, req *model.CreateAlertRuleReq) error {
	rule := &model.PromAlertRule{
		Name:        req.Name,
		PoolID:      req.PoolID,
		Expr:        req.Expr,
		Duration:    req.Duration,
		Severity:    req.Severity,
		Summary:     req.Summary,
		Description: req.Description,
		Labels:      req.Labels,
		Status:      1,
	}
	if rule.Duration == "" {
		rule.Duration = "5m"
	}
	if rule.Severity == "" {
		rule.Severity = "warning"
	}
	return s.dao.Create(ctx, rule)
}

func (s *alertRuleService) UpdateAlertRule(ctx context.Context, req *model.UpdateAlertRuleReq) error {
	rule, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Expr != "" {
		rule.Expr = req.Expr
	}
	if req.Duration != "" {
		rule.Duration = req.Duration
	}
	if req.Severity != "" {
		rule.Severity = req.Severity
	}
	if req.Summary != "" {
		rule.Summary = req.Summary
	}
	if req.Status != nil {
		rule.Status = *req.Status
	}
	return s.dao.Update(ctx, rule)
}

func (s *alertRuleService) DeleteAlertRule(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
