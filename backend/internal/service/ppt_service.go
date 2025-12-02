package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/repository"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/latex"
)

type PPTService struct {
	pptRepo         *repository.PPTRepository
	knowledgeService *KnowledgeService
	aiService       *AIService
	compiler        *latex.Compiler
	outputDir       string
}

func NewPPTService(
	pptRepo *repository.PPTRepository,
	knowledgeService *KnowledgeService,
	aiService *AIService,
	compiler *latex.Compiler,
	outputDir string,
) *PPTService {
	return &PPTService{
		pptRepo:         pptRepo,
		knowledgeService: knowledgeService,
		aiService:       aiService,
		compiler:        compiler,
		outputDir:       outputDir,
	}
}

func (s *PPTService) GeneratePPT(ctx context.Context, userID uint, title, prompt, template string, documentIDs []uint, useOpenAI bool) (*model.PPTRecord, error) {
	// Create PPT record
	ppt := &model.PPTRecord{
		UserID:   userID,
		Title:    title,
		Prompt:   prompt,
		Template: template,
		Status:   "generating",
	}

	if err := s.pptRepo.Create(ppt); err != nil {
		return nil, err
	}

	// Get context from knowledge base if document IDs provided
	var contextChunks []string
	if len(documentIDs) > 0 {
		results, err := s.knowledgeService.SearchSimilarChunks(ctx, prompt, 5)
		if err == nil {
			for _, result := range results {
				contextChunks = append(contextChunks, result.Content)
			}
		}
	}

	// Generate LaTeX content using AI
	latexContent, err := s.aiService.GenerateLaTeXPPT(ctx, prompt, contextChunks, useOpenAI)
	if err != nil {
		ppt.Status = "failed"
		ppt.ErrorMessage = err.Error()
		s.pptRepo.Update(ppt)
		return nil, err
	}

	// Extract LaTeX code from markdown code blocks if present
	latexContent = extractLatexCode(latexContent)

	ppt.LatexContent = latexContent

	// Compile LaTeX to PDF
	filename := fmt.Sprintf("ppt_%d_%d.pdf", ppt.ID, time.Now().Unix())
	pdfPath, err := s.compiler.Compile(latexContent, filename)
	if err != nil {
		ppt.Status = "failed"
		ppt.ErrorMessage = fmt.Sprintf("Compilation failed: %v", err)
		s.pptRepo.Update(ppt)
		return ppt, nil // Return ppt with latex content even if compilation fails
	}

	ppt.PDFPath = pdfPath
	ppt.Status = "completed"
	s.pptRepo.Update(ppt)

	// Create knowledge references
	for _, docID := range documentIDs {
		ref := &model.PPTKnowledgeRef{
			PPTID:      ppt.ID,
			DocumentID: docID,
		}
		s.pptRepo.CreateKnowledgeRef(ref)
	}

	return ppt, nil
}

func (s *PPTService) CompileLaTeX(pptID uint, latexContent string) error {
	ppt, err := s.pptRepo.FindByID(pptID)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("ppt_%d_%d.pdf", ppt.ID, time.Now().Unix())
	pdfPath, err := s.compiler.Compile(latexContent, filename)
	if err != nil {
		return err
	}

	ppt.LatexContent = latexContent
	ppt.PDFPath = pdfPath
	ppt.Status = "completed"
	return s.pptRepo.Update(ppt)
}

func (s *PPTService) GetPPTHistory(userID uint) ([]model.PPTRecord, error) {
	return s.pptRepo.FindByUserID(userID)
}

func (s *PPTService) GetPPT(id uint) (*model.PPTRecord, error) {
	return s.pptRepo.FindByID(id)
}

func (s *PPTService) DeletePPT(id uint) error {
	_, err := s.pptRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete PDF file if exists (simplified for now)
	
	return s.pptRepo.Delete(id)
}

func (s *PPTService) GetTemplates() []string {
	return latex.ListTemplates()
}

func extractLatexCode(content string) string {
	// Try to extract from markdown code blocks
	re := regexp.MustCompile("(?s)```latex\\s*(.+?)```")
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// If no code blocks, return as is
	return strings.TrimSpace(content)
}
