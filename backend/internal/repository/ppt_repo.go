package repository

import (
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"gorm.io/gorm"
)

type PPTRepository struct {
	db *gorm.DB
}

func NewPPTRepository(db *gorm.DB) *PPTRepository {
	return &PPTRepository{db: db}
}

func (r *PPTRepository) Create(ppt *model.PPTRecord) error {
	return r.db.Create(ppt).Error
}

func (r *PPTRepository) FindByID(id uint) (*model.PPTRecord, error) {
	var ppt model.PPTRecord
	err := r.db.First(&ppt, id).Error
	return &ppt, err
}

func (r *PPTRepository) FindByUserID(userID uint) ([]model.PPTRecord, error) {
	var ppts []model.PPTRecord
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&ppts).Error
	return ppts, err
}

func (r *PPTRepository) Update(ppt *model.PPTRecord) error {
	return r.db.Save(ppt).Error
}

func (r *PPTRepository) Delete(id uint) error {
	return r.db.Delete(&model.PPTRecord{}, id).Error
}

func (r *PPTRepository) CreateKnowledgeRef(ref *model.PPTKnowledgeRef) error {
	return r.db.Create(ref).Error
}

func (r *PPTRepository) FindKnowledgeRefsByPPTID(pptID uint) ([]model.PPTKnowledgeRef, error) {
	var refs []model.PPTKnowledgeRef
	err := r.db.Where("ppt_id = ?", pptID).Find(&refs).Error
	return refs, err
}
