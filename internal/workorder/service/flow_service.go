package service

import (
	"context"
	"errors"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"github.com/GoSimplicity/CloudOps/internal/workorder/dao"
	"go.uber.org/zap"
)

// State machine transitions
var validTransitions = map[string]map[string]string{
	"pending":     {"assign": "in_progress", "reject": "rejected"},
	"in_progress": {"resolve": "resolved", "reject": "rejected"},
	"resolved":    {"close": "closed", "reopen": "in_progress"},
	"rejected":    {"reopen": "pending"},
}

type FlowService interface {
	GetFlows(ctx context.Context, instanceID int) ([]*model.WoFlow, error)
	ExecuteAction(ctx context.Context, req *model.FlowActionReq) error
}

type flowService struct {
	flowDAO     dao.FlowDAO
	instanceDAO dao.InstanceDAO
	logger      *zap.Logger
}

func NewFlowService(flowDAO dao.FlowDAO, instanceDAO dao.InstanceDAO, logger *zap.Logger) FlowService {
	return &flowService{flowDAO: flowDAO, instanceDAO: instanceDAO, logger: logger}
}

func (s *flowService) GetFlows(ctx context.Context, instanceID int) ([]*model.WoFlow, error) {
	return s.flowDAO.ListByInstance(ctx, instanceID)
}

func (s *flowService) ExecuteAction(ctx context.Context, req *model.FlowActionReq) error {
	instance, err := s.instanceDAO.GetByID(ctx, req.InstanceID)
	if err != nil {
		return errors.New("工单不存在")
	}

	transitions, ok := validTransitions[instance.Status]
	if !ok {
		return errors.New("当前状态不支持任何操作")
	}

	toStatus, ok := transitions[req.Action]
	if !ok {
		return errors.New("当前状态不支持该操作: " + req.Action)
	}

	flow := &model.WoFlow{
		InstanceID:   req.InstanceID,
		FromStatus:   instance.Status,
		ToStatus:     toStatus,
		Action:       req.Action,
		OperatorID:   req.OperatorID,
		OperatorName: req.OperatorName,
		Remark:       req.Remark,
	}

	instance.Status = toStatus
	if req.Action == "assign" && req.AssigneeID > 0 {
		instance.AssigneeID = req.AssigneeID
	}

	if err := s.instanceDAO.Update(ctx, instance); err != nil {
		return err
	}
	return s.flowDAO.Create(ctx, flow)
}
