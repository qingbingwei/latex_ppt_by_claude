package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/api/middleware"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/service"
)

type KnowledgeHandler struct {
	knowledgeService *service.KnowledgeService
}

func NewKnowledgeHandler(knowledgeService *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{
		knowledgeService: knowledgeService,
	}
}

func (h *KnowledgeHandler) Upload(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileContent.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), file.Filename)
	
	// Read file bytes
	fileBytes := make([]byte, file.Size)
	if _, err := fileContent.Read(fileBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Save file
	filePath, err := h.knowledgeService.SaveUploadedFile(fileBytes, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create document record
	doc := &model.Document{
		UserID:   userID,
		Filename: file.Filename,
		FileType: filepath.Ext(file.Filename),
		FileSize: file.Size,
		FilePath: filePath,
		Status:   "pending",
	}

	if err := h.knowledgeService.CreateDocument(doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document record"})
		return
	}

	// Process document asynchronously
	go func() {
		h.knowledgeService.ProcessDocument(c.Request.Context(), doc, filePath)
	}()

	c.JSON(http.StatusCreated, doc)
}

func (h *KnowledgeHandler) List(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	docs, err := h.knowledgeService.GetDocumentsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get documents"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

func (h *KnowledgeHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := h.knowledgeService.GetDocument(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *KnowledgeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	if err := h.knowledgeService.DeleteDocument(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}

type SearchRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TopK == 0 {
		req.TopK = 5
	}

	results, err := h.knowledgeService.SearchSimilarChunks(c.Request.Context(), req.Query, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
