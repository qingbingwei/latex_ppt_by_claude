package model

import (
	"time"
)

type PPTRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Title        string    `gorm:"size:255" json:"title"`
	Prompt       string    `gorm:"type:text;not null" json:"prompt"`
	LatexContent string    `gorm:"type:text" json:"latex_content"`
	PDFPath      string    `gorm:"size:500" json:"pdf_path"`
	Template     string    `gorm:"size:50;default:'default'" json:"template"`
	Status       string    `gorm:"size:20;default:'pending'" json:"status"` // pending, generating, completed, failed
	ErrorMessage string    `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (PPTRecord) TableName() string {
	return "ppt_records"
}

type PPTKnowledgeRef struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PPTID      uint      `gorm:"index;not null" json:"ppt_id"`
	DocumentID uint      `gorm:"index" json:"document_id"`
	ChunkIDs   string    `gorm:"type:text" json:"chunk_ids"` // JSON string of chunk IDs
	CreatedAt  time.Time `json:"created_at"`
}

func (PPTKnowledgeRef) TableName() string {
	return "ppt_knowledge_refs"
}
