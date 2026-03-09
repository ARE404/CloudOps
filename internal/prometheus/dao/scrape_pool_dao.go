package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ScrapePoolDAO interface {
	List(ctx context.Context, req *model.ListScrapePoolsReq) ([]*model.PromScrapePool, int64, error)
	GetByID(ctx context.Context, id int) (*model.PromScrapePool, error)
	Create(ctx context.Context, pool *model.PromScrapePool) error
	Update(ctx context.Context, pool *model.PromScrapePool) error
	Delete(ctx context.Context, id int) error
}

type scrapePoolDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewScrapePoolDAO(db *gorm.DB, logger *zap.Logger) ScrapePoolDAO {
	return &scrapePoolDAO{db: db, logger: logger}
}

func (d *scrapePoolDAO) List(ctx context.Context, req *model.ListScrapePoolsReq) ([]*model.PromScrapePool, int64, error) {
	var pools []*model.PromScrapePool
	var total int64
	query := d.db.WithContext(ctx).Model(&model.PromScrapePool{}).Where("deleted_at IS NULL")
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&pools).Error; err != nil {
		return nil, 0, err
	}
	return pools, total, nil
}

func (d *scrapePoolDAO) GetByID(ctx context.Context, id int) (*model.PromScrapePool, error) {
	var pool model.PromScrapePool
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&pool).Error; err != nil {
		return nil, err
	}
	return &pool, nil
}

func (d *scrapePoolDAO) Create(ctx context.Context, pool *model.PromScrapePool) error {
	return d.db.WithContext(ctx).Create(pool).Error
}

func (d *scrapePoolDAO) Update(ctx context.Context, pool *model.PromScrapePool) error {
	return d.db.WithContext(ctx).Save(pool).Error
}

func (d *scrapePoolDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.PromScrapePool{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
