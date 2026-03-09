package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/system/dao"
	"github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx *gin.Context, req *model.LoginReq) (*model.LoginResp, error)
	Logout(ctx *gin.Context) error
	CreateUser(ctx context.Context, req *model.CreateUserReq) error
	UpdateUser(ctx context.Context, req *model.UpdateUserReq) error
	DeleteUser(ctx context.Context, id int) error
	ListUsers(ctx context.Context, req *model.ListUsersReq) (*model.PageResp, error)
	GetProfile(ctx context.Context, userID int) (*model.User, error)
}

type userService struct {
	dao        dao.UserDAO
	roleDAO    dao.RoleDAO
	jwtHandler jwt.Handler
	logger     *zap.Logger
}

func NewUserService(dao dao.UserDAO, roleDAO dao.RoleDAO, jwtHandler jwt.Handler, logger *zap.Logger) UserService {
	return &userService{dao: dao, roleDAO: roleDAO, jwtHandler: jwtHandler, logger: logger}
}

func (s *userService) Login(ctx *gin.Context, req *model.LoginReq) (*model.LoginResp, error) {
	user, err := s.dao.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}

	token, refreshToken, err := s.jwtHandler.SetLoginToken(ctx, user.ID, user.Username, user.AccountType)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}

	return &model.LoginResp{
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Username:     user.Username,
	}, nil
}

func (s *userService) Logout(ctx *gin.Context) error {
	return s.jwtHandler.ClearToken(ctx)
}

func (s *userService) CreateUser(ctx context.Context, req *model.CreateUserReq) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user := &model.User{
		Username: req.Username,
		Password: string(hash),
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1,
	}

	if len(req.RoleIDs) > 0 {
		roles, err := s.roleDAO.GetByIDs(ctx, req.RoleIDs)
		if err != nil {
			return err
		}
		user.Roles = roles
	}

	return s.dao.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, req *model.UpdateUserReq) error {
	user, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if len(req.RoleIDs) > 0 {
		roles, err := s.roleDAO.GetByIDs(ctx, req.RoleIDs)
		if err != nil {
			return err
		}
		user.Roles = roles
	}

	return s.dao.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, req *model.ListUsersReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*model.User, error) {
	return s.dao.GetByID(ctx, userID)
}
