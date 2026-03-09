package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/dao"
	"go.uber.org/zap"
)

type MenuService interface {
	ListMenus(ctx context.Context) ([]*model.Menu, error)
	CreateMenu(ctx context.Context, req *model.CreateMenuReq) error
	UpdateMenu(ctx context.Context, req *model.UpdateMenuReq) error
	DeleteMenu(ctx context.Context, id int) error
}

type menuService struct {
	dao    dao.MenuDAO
	logger *zap.Logger
}

func NewMenuService(dao dao.MenuDAO, logger *zap.Logger) MenuService {
	return &menuService{dao: dao, logger: logger}
}

func (s *menuService) ListMenus(ctx context.Context) ([]*model.Menu, error) {
	return s.dao.List(ctx)
}

func (s *menuService) CreateMenu(ctx context.Context, req *model.CreateMenuReq) error {
	menu := &model.Menu{
		Name:     req.Name,
		Path:     req.Path,
		Icon:     req.Icon,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		MenuType: req.MenuType,
	}
	return s.dao.Create(ctx, menu)
}

func (s *menuService) UpdateMenu(ctx context.Context, req *model.UpdateMenuReq) error {
	menu, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	if req.Hidden != nil {
		menu.Hidden = *req.Hidden
	}
	return s.dao.Update(ctx, menu)
}

func (s *menuService) DeleteMenu(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
