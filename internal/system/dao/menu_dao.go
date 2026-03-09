package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MenuDAO interface {
	List(ctx context.Context) ([]*model.Menu, error)
	GetByID(ctx context.Context, id int) (*model.Menu, error)
	Create(ctx context.Context, menu *model.Menu) error
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id int) error
	GetByIDs(ctx context.Context, ids []int) ([]*model.Menu, error)
}

type menuDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewMenuDAO(db *gorm.DB, logger *zap.Logger) MenuDAO {
	return &menuDAO{db: db, logger: logger}
}

func (d *menuDAO) List(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	if err := d.db.WithContext(ctx).Where("deleted_at IS NULL").Order("sort ASC").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (d *menuDAO) GetByID(ctx context.Context, id int) (*model.Menu, error) {
	var menu model.Menu
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (d *menuDAO) Create(ctx context.Context, menu *model.Menu) error {
	return d.db.WithContext(ctx).Create(menu).Error
}

func (d *menuDAO) Update(ctx context.Context, menu *model.Menu) error {
	return d.db.WithContext(ctx).Save(menu).Error
}

func (d *menuDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.Menu{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (d *menuDAO) GetByIDs(ctx context.Context, ids []int) ([]*model.Menu, error) {
	var menus []*model.Menu
	if err := d.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}
