package dao

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserDAO interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	List(ctx context.Context, req *model.ListUsersReq) ([]*model.User, int64, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int) error
}

type userDAO struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserDAO(db *gorm.DB, logger *zap.Logger) UserDAO {
	return &userDAO{db: db, logger: logger}
}

func (d *userDAO) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).Where("username = ? AND deleted_at IS NULL", username).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	err := d.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAO) List(ctx context.Context, req *model.ListUsersReq) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	query := d.db.WithContext(ctx).Model(&model.User{}).Where("deleted_at IS NULL")
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Preload("Roles").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (d *userDAO) Create(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userDAO) Update(ctx context.Context, user *model.User) error {
	return d.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (d *userDAO) Delete(ctx context.Context, id int) error {
	return d.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
