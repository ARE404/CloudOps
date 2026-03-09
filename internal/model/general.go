package model

import "time"

// Model 基础模型，包含软删除
type Model struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// PageReq 分页请求
type PageReq struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *PageReq) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	return (p.Page - 1) * p.PageSize
}

func (p *PageReq) GetLimit() int {
	if p.PageSize <= 0 {
		return 20
	}
	return p.PageSize
}

// PageResp 分页响应
type PageResp struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
