package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FlowDAO interface {
	ListByInstance(ctx context.Context, instanceID int) ([]*model.WoFlow, error)
	Create(ctx context.Context, flow *model.WoFlow) error
}

type flowDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewFlowDAO(db *gorm.DB, logger *zap.Logger) FlowDAO {
	return &flowDAO{db: db, logger: logger}
}

func (d *flowDAO) ListByInstance(ctx context.Context, instanceID int) ([]*model.WoFlow, error) {
	var flows []*model.WoFlow
	if err := d.db.WithContext(ctx).Where("instance_id = ?", instanceID).Order("created_at ASC").Find(&flows).Error; err != nil {
		return nil, err
	}
	return flows, nil
}

func (d *flowDAO) Create(ctx context.Context, flow *model.WoFlow) error {
	return d.db.WithContext(ctx).Create(flow).Error
}
