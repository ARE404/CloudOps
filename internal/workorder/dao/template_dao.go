package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TemplateDAO interface {
	List(ctx context.Context, req *model.ListTemplatesReq) ([]*model.WoTemplate, int64, error)
	GetByID(ctx context.Context, id int) (*model.WoTemplate, error)
	Create(ctx context.Context, t *model.WoTemplate) error
	Update(ctx context.Context, t *model.WoTemplate) error
	Delete(ctx context.Context, id int) error
}

type templateDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewTemplateDAO(db *gorm.DB, logger *zap.Logger) TemplateDAO {
	return &templateDAO{db: db, logger: logger}
}

func (d *templateDAO) List(ctx context.Context, req *model.ListTemplatesReq) ([]*model.WoTemplate, int64, error) {
	var templates []*model.WoTemplate
	var total int64
	query := d.db.WithContext(ctx).Model(&model.WoTemplate{}).Where("deleted_at IS NULL")
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&templates).Error; err != nil {
		return nil, 0, err
	}
	return templates, total, nil
}

func (d *templateDAO) GetByID(ctx context.Context, id int) (*model.WoTemplate, error) {
	var t model.WoTemplate
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (d *templateDAO) Create(ctx context.Context, t *model.WoTemplate) error {
	return d.db.WithContext(ctx).Create(t).Error
}

func (d *templateDAO) Update(ctx context.Context, t *model.WoTemplate) error {
	return d.db.WithContext(ctx).Save(t).Error
}

func (d *templateDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.WoTemplate{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
