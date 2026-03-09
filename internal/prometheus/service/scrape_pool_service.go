package service

import (
	"context"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/dao"
	"go.uber.org/zap"
)

type ScrapePoolService interface {
	ListScrapePools(ctx context.Context, req *model.ListScrapePoolsReq) (*model.PageResp, error)
	CreateScrapePool(ctx context.Context, req *model.CreateScrapePoolReq) error
	UpdateScrapePool(ctx context.Context, req *model.UpdateScrapePoolReq) error
	DeleteScrapePool(ctx context.Context, id int) error
}

type scrapePoolService struct {
	dao    dao.ScrapePoolDAO
	logger *zap.Logger
}

func NewScrapePoolService(dao dao.ScrapePoolDAO, logger *zap.Logger) ScrapePoolService {
	return &scrapePoolService{dao: dao, logger: logger}
}

func (s *scrapePoolService) ListScrapePools(ctx context.Context, req *model.ListScrapePoolsReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *scrapePoolService) CreateScrapePool(ctx context.Context, req *model.CreateScrapePoolReq) error {
	pool := &model.PromScrapePool{
		Name:             req.Name,
		PrometheusAddr:   req.PrometheusAddr,
		AlertManagerAddr: req.AlertManagerAddr,
		ScrapeInterval:   req.ScrapeInterval,
		Description:      req.Description,
	}
	if pool.ScrapeInterval == 0 {
		pool.ScrapeInterval = 15
	}
	return s.dao.Create(ctx, pool)
}

func (s *scrapePoolService) UpdateScrapePool(ctx context.Context, req *model.UpdateScrapePoolReq) error {
	pool, err := s.dao.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if req.PrometheusAddr != "" {
		pool.PrometheusAddr = req.PrometheusAddr
	}
	if req.AlertManagerAddr != "" {
		pool.AlertManagerAddr = req.AlertManagerAddr
	}
	if req.ScrapeInterval > 0 {
		pool.ScrapeInterval = req.ScrapeInterval
	}
	if req.Description != "" {
		pool.Description = req.Description
	}
	return s.dao.Update(ctx, pool)
}

func (s *scrapePoolService) DeleteScrapePool(ctx context.Context, id int) error {
	return s.dao.Delete(ctx, id)
}
