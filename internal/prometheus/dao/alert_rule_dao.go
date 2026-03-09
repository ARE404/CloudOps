package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlertRuleDAO interface {
	List(ctx context.Context, req *model.ListAlertRulesReq) ([]*model.PromAlertRule, int64, error)
	GetByID(ctx context.Context, id int) (*model.PromAlertRule, error)
	Create(ctx context.Context, rule *model.PromAlertRule) error
	Update(ctx context.Context, rule *model.PromAlertRule) error
	Delete(ctx context.Context, id int) error
}

type alertRuleDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAlertRuleDAO(db *gorm.DB, logger *zap.Logger) AlertRuleDAO {
	return &alertRuleDAO{db: db, logger: logger}
}

func (d *alertRuleDAO) List(ctx context.Context, req *model.ListAlertRulesReq) ([]*model.PromAlertRule, int64, error) {
	var rules []*model.PromAlertRule
	var total int64
	query := d.db.WithContext(ctx).Model(&model.PromAlertRule{}).Where("deleted_at IS NULL")
	if req.PoolID > 0 {
		query = query.Where("pool_id = ?", req.PoolID)
	}
	if req.Severity != "" {
		query = query.Where("severity = ?", req.Severity)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

func (d *alertRuleDAO) GetByID(ctx context.Context, id int) (*model.PromAlertRule, error) {
	var rule model.PromAlertRule
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&rule).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

func (d *alertRuleDAO) Create(ctx context.Context, rule *model.PromAlertRule) error {
	return d.db.WithContext(ctx).Create(rule).Error
}

func (d *alertRuleDAO) Update(ctx context.Context, rule *model.PromAlertRule) error {
	return d.db.WithContext(ctx).Save(rule).Error
}

func (d *alertRuleDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.PromAlertRule{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
