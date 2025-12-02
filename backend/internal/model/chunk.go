package model

import (
	"time"
)

type Chunk struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DocumentID uint      `gorm:"index;not null" json:"document_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	ChunkIndex int       `json:"chunk_index"`
	VectorID   string    `gorm:"size:100" json:"vector_id"`
	Metadata   string    `gorm:"type:text" json:"metadata"` // JSON string
	CreatedAt  time.Time `json:"created_at"`
}

func (Chunk) TableName() string {
	return "chunks"
}
