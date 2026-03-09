package model

import "time"

// WoTemplate 工单模板
type WoTemplate struct {
	Model
	Name        string `json:"name" gorm:"size:64;not null"`
	Description string `json:"description" gorm:"size:255"`
	FormSchema  string `json:"form_schema" gorm:"type:text"` // JSON
	Category    string `json:"category" gorm:"size:32"`
	Status      int8   `json:"status" gorm:"default:1"`
	CreateUserID int   `json:"create_user_id"`
}

// WoInstance 工单实例
type WoInstance struct {
	Model
	Title       string     `json:"title" gorm:"size:128;not null"`
	TemplateID  int        `json:"template_id" gorm:"not null;index"`
	FormData    string     `json:"form_data" gorm:"type:text"` // JSON
	Status      string     `json:"status" gorm:"size:16;default:'pending'"` // pending|in_progress|resolved|closed|rejected
	Priority    int8       `json:"priority" gorm:"default:2"` // 1=low,2=medium,3=high
	AssigneeID  int        `json:"assignee_id" gorm:"index"`
	ReporterID  int        `json:"reporter_id" gorm:"index"`
	ReporterName string    `json:"reporter_name" gorm:"size:64"`
	DueAt       *time.Time `json:"due_at"`
}

// WoFlow 工单流转记录
type WoFlow struct {
	Model
	InstanceID int    `json:"instance_id" gorm:"not null;index"`
	FromStatus string `json:"from_status" gorm:"size:16"`
	ToStatus   string `json:"to_status" gorm:"size:16;not null"`
	Action     string `json:"action" gorm:"size:32"` // submit|assign|resolve|close|reject|reopen
	OperatorID int    `json:"operator_id"`
	OperatorName string `json:"operator_name" gorm:"size:64"`
	Remark     string `json:"remark" gorm:"size:255"`
}

// WoComment 工单评论
type WoComment struct {
	Model
	InstanceID int    `json:"instance_id" gorm:"not null;index"`
	UserID     int    `json:"user_id" gorm:"not null"`
	Username   string `json:"username" gorm:"size:64"`
	Content    string `json:"content" gorm:"type:text;not null"`
}

// --- Requests ---

type ListTemplatesReq struct {
	PageReq
	Category string `json:"category" form:"category"`
}

type CreateTemplateReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	FormSchema  string `json:"form_schema"`
	Category    string `json:"category"`
}

type UpdateTemplateReq struct {
	ID          int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FormSchema  string `json:"form_schema"`
	Status      *int8  `json:"status"`
}

type ListInstancesReq struct {
	PageReq
	Status     string `json:"status" form:"status"`
	AssigneeID int    `json:"assignee_id" form:"assignee_id"`
	ReporterID int    `json:"reporter_id" form:"reporter_id"`
}

type CreateInstanceReq struct {
	Title      string     `json:"title" binding:"required"`
	TemplateID int        `json:"template_id" binding:"required"`
	FormData   string     `json:"form_data"`
	Priority   int8       `json:"priority"`
	AssigneeID int        `json:"assignee_id"`
	DueAt      *time.Time `json:"due_at"`
	ReporterID int        `json:"-"`
	ReporterName string   `json:"-"`
}

type FlowActionReq struct {
	InstanceID   int    `json:"-"`
	Action       string `json:"action" binding:"required"` // assign|resolve|close|reject|reopen
	AssigneeID   int    `json:"assignee_id"`
	Remark       string `json:"remark"`
	OperatorID   int    `json:"-"`
	OperatorName string `json:"-"`
}

type CreateCommentReq struct {
	InstanceID int    `json:"-"`
	Content    string `json:"content" binding:"required"`
	UserID     int    `json:"-"`
	Username   string `json:"-"`
}

type ListCommentsReq struct {
	InstanceID int `json:"-"`
}
