package model

// AiKnowledgeChunk AI知识库文档块
type AiKnowledgeChunk struct {
	Model
	Title     string  `json:"title" gorm:"size:255;not null"`
	Content   string  `json:"content" gorm:"type:text;not null"`
	Source    string  `json:"source" gorm:"size:255"`
	Category  string  `json:"category" gorm:"size:64"`
	Embedding string  `json:"embedding" gorm:"type:longtext"` // JSON float64 array
	Score     float64 `json:"score,omitempty" gorm:"-"`
}

// --- Requests ---

type ChatReq struct {
	Message   string `json:"message" binding:"required"`
	SessionID string `json:"session_id"`
	UseRAG    bool   `json:"use_rag"`
}

type ChatResp struct {
	Answer    string `json:"answer"`
	SessionID string `json:"session_id"`
}

type AddKnowledgeReq struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Source   string `json:"source"`
	Category string `json:"category"`
}

type ListKnowledgeReq struct {
	PageReq
	Category string `json:"category" form:"category"`
}
