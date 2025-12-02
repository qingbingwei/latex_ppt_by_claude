package service

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/repository"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/embedding"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/parser"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/vectordb"
)

type KnowledgeService struct {
	docRepo         *repository.DocumentRepository
	embeddingClient *embedding.OpenAIEmbedding
	vectorDB        *vectordb.MilvusClient
	uploadDir       string
}

func NewKnowledgeService(
	docRepo *repository.DocumentRepository,
	embeddingClient *embedding.OpenAIEmbedding,
	vectorDB *vectordb.MilvusClient,
	uploadDir string,
) *KnowledgeService {
	return &KnowledgeService{
		docRepo:         docRepo,
		embeddingClient: embeddingClient,
		vectorDB:        vectorDB,
		uploadDir:       uploadDir,
	}
}

func (s *KnowledgeService) ProcessDocument(ctx context.Context, doc *model.Document, filePath string) error {
	log.Printf("Starting to process document: %s (ID: %d)", doc.Filename, doc.ID)

	// Update status to processing
	doc.Status = "processing"
	if err := s.docRepo.Update(doc); err != nil {
		log.Printf("Failed to update document status to processing: %v", err)
		return err
	}

	// Parse document
	p, err := parser.GetParser(doc.Filename)
	if err != nil {
		log.Printf("Failed to get parser for %s: %v", doc.Filename, err)
		doc.Status = "failed"
		s.docRepo.Update(doc)
		return err
	}

	content, err := p.Parse(filePath)
	if err != nil {
		log.Printf("Failed to parse document %s: %v", doc.Filename, err)
		doc.Status = "failed"
		s.docRepo.Update(doc)
		return err
	}

	log.Printf("Parsed document %s, content length: %d chars", doc.Filename, len(content))

	// Split into chunks
	chunks := parser.SplitIntoChunks(content, 200, 50)
	log.Printf("Split document %s into %d chunks", doc.Filename, len(chunks))

	successCount := 0
	// Generate embeddings and store in vector DB
	for i, chunkText := range chunks {
		// Create chunk record
		chunk := &model.Chunk{
			DocumentID: doc.ID,
			Content:    chunkText,
			ChunkIndex: i,
		}

		if err := s.docRepo.CreateChunk(chunk); err != nil {
			log.Printf("Failed to create chunk %d: %v", i, err)
			continue
		}

		// Generate embedding - 使用新的 context 避免请求结束后 context 被取消
		emb, err := s.embeddingClient.GenerateEmbedding(context.Background(), chunkText)
		if err != nil {
			log.Printf("Failed to generate embedding for chunk %d: %v", i, err)
			continue
		}

		log.Printf("Generated embedding for chunk %d, vector length: %d", i, len(emb))

		// Store in Milvus
		vectorID, err := s.vectorDB.Insert(context.Background(), int64(chunk.ID), int64(doc.ID), chunkText, emb)
		if err != nil {
			log.Printf("Failed to insert chunk %d into Milvus: %v", i, err)
			continue
		}

		// Update chunk with vector ID
		chunk.VectorID = vectorID
		successCount++
	}

	// Update document status
	doc.Status = "completed"
	doc.ChunkCount = successCount
	log.Printf("Document %s processing completed, %d/%d chunks successful", doc.Filename, successCount, len(chunks))
	return s.docRepo.Update(doc)
}

func (s *KnowledgeService) SaveUploadedFile(file []byte, filename string) (string, error) {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(s.uploadDir, filename)
	return filePath, os.WriteFile(filePath, file, 0644)
}

func (s *KnowledgeService) SearchSimilarChunks(ctx context.Context, query string, topK int) ([]vectordb.SearchResult, error) {
	// Generate query embedding
	emb, err := s.embeddingClient.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	// Search in Milvus
	return s.vectorDB.Search(ctx, emb, topK)
}

func (s *KnowledgeService) GetDocumentsByUser(userID uint) ([]model.Document, error) {
	return s.docRepo.FindByUserID(userID)
}

func (s *KnowledgeService) GetDocument(id uint) (*model.Document, error) {
	return s.docRepo.FindByID(id)
}

func (s *KnowledgeService) DeleteDocument(id uint) error {
	doc, err := s.docRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete file
	if doc.FilePath != "" {
		os.Remove(doc.FilePath)
	}

	// Delete from database (chunks will be cascade deleted)
	return s.docRepo.Delete(id)
}

func (s *KnowledgeService) CreateDocument(doc *model.Document) error {
	return s.docRepo.Create(doc)
}
