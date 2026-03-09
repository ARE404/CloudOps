package model

// PromScrapePool Prometheus 抓取池（对应 Prometheus 实例）
type PromScrapePool struct {
	Model
	Name               string `json:"name" gorm:"uniqueIndex;size:64;not null"`
	PrometheusAddr     string `json:"prometheus_addr" gorm:"size:255"`
	AlertManagerAddr   string `json:"alert_manager_addr" gorm:"size:255"`
	ExternalLabels     string `json:"external_labels" gorm:"type:text"` // JSON
	ScrapeInterval     int    `json:"scrape_interval" gorm:"default:15"`
	ScrapeTimeout      int    `json:"scrape_timeout" gorm:"default:10"`
	RemoteWriteAddr    string `json:"remote_write_addr" gorm:"size:255"`
	Description        string `json:"description" gorm:"size:255"`
	CreateUserID       int    `json:"create_user_id"`
}

// PromAlertRule 告警规则
type PromAlertRule struct {
	Model
	Name        string `json:"name" gorm:"size:128;not null"`
	PoolID      int    `json:"pool_id" gorm:"not null;index"`
	Expr        string `json:"expr" gorm:"type:text;not null"`
	Duration    string `json:"duration" gorm:"size:32;default:'5m'"`
	Severity    string `json:"severity" gorm:"size:16;default:'warning'"`
	Summary     string `json:"summary" gorm:"size:255"`
	Description string `json:"description" gorm:"type:text"`
	Labels      string `json:"labels" gorm:"type:text"` // JSON key=value
	Status      int8   `json:"status" gorm:"default:1"`
}

// PromAlertEvent 告警事件（AlertManager 推送）
type PromAlertEvent struct {
	Model
	AlertName    string `json:"alert_name" gorm:"size:128;not null"`
	PoolID       int    `json:"pool_id" gorm:"index"`
	RuleID       int    `json:"rule_id" gorm:"index"`
	Fingerprint  string `json:"fingerprint" gorm:"size:64;index"`
	Status       string `json:"status" gorm:"size:16;default:'firing'"` // firing | resolved
	Severity     string `json:"severity" gorm:"size:16"`
	Summary      string `json:"summary" gorm:"size:255"`
	Description  string `json:"description" gorm:"type:text"`
	Labels       string `json:"labels" gorm:"type:text"`  // JSON
	SendGroupID  int    `json:"send_group_id" gorm:"index"`
	NotifiedAt   string `json:"notified_at" gorm:"size:32"`
}

// PromSendGroup 告警发送组
type PromSendGroup struct {
	Model
	Name        string `json:"name" gorm:"size:64;not null"`
	PoolID      int    `json:"pool_id" gorm:"index"`
	Channels    string `json:"channels" gorm:"type:text"` // JSON: [{type,target}]
	RepeatInterval int `json:"repeat_interval" gorm:"default:3600"`
	Description string `json:"description" gorm:"size:255"`
}

// --- Requests ---

type ListScrapePoolsReq struct {
	PageReq
	Name string `json:"name" form:"name"`
}

type CreateScrapePoolReq struct {
	Name             string `json:"name" binding:"required"`
	PrometheusAddr   string `json:"prometheus_addr" binding:"required"`
	AlertManagerAddr string `json:"alert_manager_addr"`
	ScrapeInterval   int    `json:"scrape_interval"`
	Description      string `json:"description"`
}

type UpdateScrapePoolReq struct {
	ID               int    `json:"-"`
	PrometheusAddr   string `json:"prometheus_addr"`
	AlertManagerAddr string `json:"alert_manager_addr"`
	ScrapeInterval   int    `json:"scrape_interval"`
	Description      string `json:"description"`
}

type ListAlertRulesReq struct {
	PageReq
	PoolID   int    `json:"pool_id" form:"pool_id"`
	Severity string `json:"severity" form:"severity"`
}

type CreateAlertRuleReq struct {
	Name        string `json:"name" binding:"required"`
	PoolID      int    `json:"pool_id" binding:"required"`
	Expr        string `json:"expr" binding:"required"`
	Duration    string `json:"duration"`
	Severity    string `json:"severity"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Labels      string `json:"labels"`
}

type UpdateAlertRuleReq struct {
	ID          int    `json:"-"`
	Expr        string `json:"expr"`
	Duration    string `json:"duration"`
	Severity    string `json:"severity"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Status      *int8  `json:"status"`
}

type ListAlertEventsReq struct {
	PageReq
	PoolID  int    `json:"pool_id" form:"pool_id"`
	Status  string `json:"status" form:"status"`
}

type ListSendGroupsReq struct {
	PageReq
	PoolID int `json:"pool_id" form:"pool_id"`
}

type CreateSendGroupReq struct {
	Name           string `json:"name" binding:"required"`
	PoolID         int    `json:"pool_id" binding:"required"`
	Channels       string `json:"channels"`
	RepeatInterval int    `json:"repeat_interval"`
	Description    string `json:"description"`
}

type UpdateSendGroupReq struct {
	ID             int    `json:"-"`
	Channels       string `json:"channels"`
	RepeatInterval int    `json:"repeat_interval"`
	Description    string `json:"description"`
}
