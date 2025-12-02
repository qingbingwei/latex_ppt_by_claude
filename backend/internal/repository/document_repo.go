package repository

import (
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"gorm.io/gorm"
)

type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) Create(doc *model.Document) error {
	return r.db.Create(doc).Error
}

func (r *DocumentRepository) FindByID(id uint) (*model.Document, error) {
	var doc model.Document
	err := r.db.First(&doc, id).Error
	return &doc, err
}

func (r *DocumentRepository) FindByUserID(userID uint) ([]model.Document, error) {
	var docs []model.Document
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&docs).Error
	return docs, err
}

func (r *DocumentRepository) Update(doc *model.Document) error {
	return r.db.Save(doc).Error
}

func (r *DocumentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Document{}, id).Error
}

func (r *DocumentRepository) CreateChunk(chunk *model.Chunk) error {
	return r.db.Create(chunk).Error
}

func (r *DocumentRepository) FindChunksByDocumentID(docID uint) ([]model.Chunk, error) {
	var chunks []model.Chunk
	err := r.db.Where("document_id = ?", docID).Order("chunk_index").Find(&chunks).Error
	return chunks, err
}

func (r *DocumentRepository) FindChunkByID(id uint) (*model.Chunk, error) {
	var chunk model.Chunk
	err := r.db.First(&chunk, id).Error
	return &chunk, err
}
