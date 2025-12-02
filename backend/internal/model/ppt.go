package model

import (
	"database/sql/driver"
	"encoding/json"
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

type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

type PPTKnowledgeRef struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PPTID      uint      `gorm:"index;not null" json:"ppt_id"`
	DocumentID uint      `gorm:"index" json:"document_id"`
	ChunkIDs   IntArray  `gorm:"type:jsonb" json:"chunk_ids"`
	CreatedAt  time.Time `json:"created_at"`
}

func (PPTKnowledgeRef) TableName() string {
	return "ppt_knowledge_refs"
}
