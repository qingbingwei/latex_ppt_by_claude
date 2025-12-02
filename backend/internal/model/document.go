package model

import (
	"time"
)

type Document struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	Filename   string    `gorm:"size:255;not null" json:"filename"`
	FileType   string    `gorm:"size:50;not null" json:"file_type"`
	FileSize   int64     `json:"file_size"`
	FilePath   string    `gorm:"size:500" json:"file_path"`
	Status     string    `gorm:"size:20;default:'pending'" json:"status"` // pending, processing, completed, failed
	ChunkCount int       `gorm:"default:0" json:"chunk_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Document) TableName() string {
	return "documents"
}
