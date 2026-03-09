package service

import (
	"context"
	"fmt"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/dao"
	"go.uber.org/zap"
)

type RoleService interface {
	ListRoles(ctx context.Context, req *model.ListRolesReq) (*model.PageResp, error)
	CreateRole(ctx context.Context, req *model.CreateRoleReq) error
	UpdateRole(ctx context.Context, req *model.UpdateRoleReq) error
	DeleteRole(ctx context.Context, id int) error
}

type roleService struct {
	dao     dao.RoleDAO
	menuDAO dao.MenuDAO
	logger  *zap.Logger
}

func NewRoleService(dao dao.RoleDAO, menuDAO dao.MenuDAO, logger *zap.Logger) RoleService {
	return &roleService{dao: dao, menuDAO: menuDAO, logger: logger}
}

func (s *roleService) ListRoles(ctx context.Context, req *model.ListRolesReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *roleService) CreateRole(ctx context.Context, req *model.CreateRoleReq) error {
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1,
	}
	if len(req.MenuIDs) > 0 {
		menus, err := s.menuDAO.GetByIDs(ctx, req.MenuIDs)
		if err != nil {
			return err
		}
		role.Menus = menus
	}
	return s.dao.Create(ctx, role)
}

func (s *roleService) UpdateRole(ctx context.Context, req *model.UpdateRoleReq) error {
	role, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}
	if len(req.MenuIDs) > 0 {
		menus, err := s.menuDAO.GetByIDs(ctx, req.MenuIDs)
		if err != nil {
			return err
		}
		role.Menus = menus
	}
	return s.dao.Update(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
