package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SendGroupDAO interface {
	List(ctx context.Context, req *model.ListSendGroupsReq) ([]*model.PromSendGroup, int64, error)
	GetByID(ctx context.Context, id int) (*model.PromSendGroup, error)
	Create(ctx context.Context, group *model.PromSendGroup) error
	Update(ctx context.Context, group *model.PromSendGroup) error
	Delete(ctx context.Context, id int) error
}

type sendGroupDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewSendGroupDAO(db *gorm.DB, logger *zap.Logger) SendGroupDAO {
	return &sendGroupDAO{db: db, logger: logger}
}

func (d *sendGroupDAO) List(ctx context.Context, req *model.ListSendGroupsReq) ([]*model.PromSendGroup, int64, error) {
	var groups []*model.PromSendGroup
	var total int64
	query := d.db.WithContext(ctx).Model(&model.PromSendGroup{}).Where("deleted_at IS NULL")
	if req.PoolID > 0 {
		query = query.Where("pool_id = ?", req.PoolID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&groups).Error; err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}

func (d *sendGroupDAO) GetByID(ctx context.Context, id int) (*model.PromSendGroup, error) {
	var group model.PromSendGroup
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (d *sendGroupDAO) Create(ctx context.Context, group *model.PromSendGroup) error {
	return d.db.WithContext(ctx).Create(group).Error
}

func (d *sendGroupDAO) Update(ctx context.Context, group *model.PromSendGroup) error {
	return d.db.WithContext(ctx).Save(group).Error
}

func (d *sendGroupDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.PromSendGroup{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
