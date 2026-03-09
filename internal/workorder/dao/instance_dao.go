package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InstanceDAO interface {
	List(ctx context.Context, req *model.ListInstancesReq) ([]*model.WoInstance, int64, error)
	GetByID(ctx context.Context, id int) (*model.WoInstance, error)
	Create(ctx context.Context, instance *model.WoInstance) error
	Update(ctx context.Context, instance *model.WoInstance) error
}

type instanceDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewInstanceDAO(db *gorm.DB, logger *zap.Logger) InstanceDAO {
	return &instanceDAO{db: db, logger: logger}
}

func (d *instanceDAO) List(ctx context.Context, req *model.ListInstancesReq) ([]*model.WoInstance, int64, error) {
	var instances []*model.WoInstance
	var total int64
	query := d.db.WithContext(ctx).Model(&model.WoInstance{}).Where("deleted_at IS NULL")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.AssigneeID > 0 {
		query = query.Where("assignee_id = ?", req.AssigneeID)
	}
	if req.ReporterID > 0 {
		query = query.Where("reporter_id = ?", req.ReporterID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&instances).Error; err != nil {
		return nil, 0, err
	}
	return instances, total, nil
}

func (d *instanceDAO) GetByID(ctx context.Context, id int) (*model.WoInstance, error) {
	var instance model.WoInstance
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&instance).Error; err != nil {
		return nil, err
	}
	return &instance, nil
}

func (d *instanceDAO) Create(ctx context.Context, instance *model.WoInstance) error {
	return d.db.WithContext(ctx).Create(instance).Error
}

func (d *instanceDAO) Update(ctx context.Context, instance *model.WoInstance) error {
	return d.db.WithContext(ctx).Save(instance).Error
}
