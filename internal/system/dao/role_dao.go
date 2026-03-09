package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleDAO interface {
	List(ctx context.Context, req *model.ListRolesReq) ([]*model.Role, int64, error)
	GetByID(ctx context.Context, id int) (*model.Role, error)
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id int) error
	GetByIDs(ctx context.Context, ids []int) ([]*model.Role, error)
}

type roleDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRoleDAO(db *gorm.DB, logger *zap.Logger) RoleDAO {
	return &roleDAO{db: db, logger: logger}
}

func (d *roleDAO) List(ctx context.Context, req *model.ListRolesReq) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64
	query := d.db.WithContext(ctx).Model(&model.Role{}).Where("deleted_at IS NULL")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("Menus").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func (d *roleDAO) GetByID(ctx context.Context, id int) (*model.Role, error) {
	var role model.Role
	if err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Preload("Menus").First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (d *roleDAO) Create(ctx context.Context, role *model.Role) error {
	return d.db.WithContext(ctx).Create(role).Error
}

func (d *roleDAO) Update(ctx context.Context, role *model.Role) error {
	return d.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(role).Error
}

func (d *roleDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.Role{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (d *roleDAO) GetByIDs(ctx context.Context, ids []int) ([]*model.Role, error) {
	var roles []*model.Role
	if err := d.db.WithContext(ctx).Where("id IN ? AND deleted_at IS NULL", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
