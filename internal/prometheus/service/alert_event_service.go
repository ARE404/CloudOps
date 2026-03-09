package service

import (
	"context"
	"encoding/json"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/dao"
	"github.com/GoSimplicity/CloudOps/internal/prometheus/notify"
	"go.uber.org/zap"
)

type AlertEventService interface {
	ListAlertEvents(ctx context.Context, req *model.ListAlertEventsReq) (*model.PageResp, error)
	ReceiveAlert(ctx context.Context, event *model.PromAlertEvent) error
}

type alertEventService struct {
	dao          dao.AlertEventDAO
	sendGroupDAO dao.SendGroupDAO
	notifiers    []notify.Notifier
	logger       *zap.Logger
}

func NewAlertEventService(
	dao dao.AlertEventDAO,
	sendGroupDAO dao.SendGroupDAO,
	logger *zap.Logger,
) AlertEventService {
	return &alertEventService{
		dao:          dao,
		sendGroupDAO: sendGroupDAO,
		notifiers: []notify.Notifier{
			notify.NewDingTalkNotifier(),
			notify.NewEmailNotifier(),
			notify.NewWebhookNotifier(),
		},
		logger: logger,
	}
}

func (s *alertEventService) ListAlertEvents(ctx context.Context, req *model.ListAlertEventsReq) (*model.PageResp, error) {
	list, total, err := s.dao.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return &model.PageResp{Total: total, List: list}, nil
}

func (s *alertEventService) ReceiveAlert(ctx context.Context, event *model.PromAlertEvent) error {
	if err := s.dao.Create(ctx, event); err != nil {
		return err
	}

	// async notify
	go func() {
		sendGroup, err := s.sendGroupDAO.GetByID(context.Background(), event.SendGroupID)
		if err != nil || sendGroup == nil {
			return
		}
		var channels []notify.NotifyChannel
		if err := json.Unmarshal([]byte(sendGroup.Channels), &channels); err != nil {
			s.logger.Error("解析通知渠道失败", zap.Error(err))
			return
		}
		for _, n := range s.notifiers {
			if err := n.Send(context.Background(), event.AlertName, event.Description, channels); err != nil {
				s.logger.Error("发送通知失败", zap.Error(err))
			}
		}
	}()

	return nil
}
