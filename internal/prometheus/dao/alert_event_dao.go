package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlertEventDAO interface {
	List(ctx context.Context, req *model.ListAlertEventsReq) ([]*model.PromAlertEvent, int64, error)
	Create(ctx context.Context, event *model.PromAlertEvent) error
	UpdateStatus(ctx context.Context, fingerprint string, status string) error
}

type alertEventDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAlertEventDAO(db *gorm.DB, logger *zap.Logger) AlertEventDAO {
	return &alertEventDAO{db: db, logger: logger}
}

func (d *alertEventDAO) List(ctx context.Context, req *model.ListAlertEventsReq) ([]*model.PromAlertEvent, int64, error) {
	var events []*model.PromAlertEvent
	var total int64
	query := d.db.WithContext(ctx).Model(&model.PromAlertEvent{}).Where("deleted_at IS NULL")
	if req.PoolID > 0 {
		query = query.Where("pool_id = ?", req.PoolID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&events).Error; err != nil {
		return nil, 0, err
	}
	return events, total, nil
}

func (d *alertEventDAO) Create(ctx context.Context, event *model.PromAlertEvent) error {
	return d.db.WithContext(ctx).Create(event).Error
}

func (d *alertEventDAO) UpdateStatus(ctx context.Context, fingerprint string, status string) error {
	return d.db.WithContext(ctx).Model(&model.PromAlertEvent{}).
		Where("fingerprint = ?", fingerprint).
		Update("status", status).Error
}
