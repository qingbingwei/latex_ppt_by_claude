package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/api/middleware"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/service"
)

type PPTHandler struct {
	pptService *service.PPTService
}

func NewPPTHandler(pptService *service.PPTService) *PPTHandler {
	return &PPTHandler{
		pptService: pptService,
	}
}

type GenerateRequest struct {
	Title       string `json:"title" binding:"required"`
	Prompt      string `json:"prompt" binding:"required"`
	Template    string `json:"template"`
	DocumentIDs []uint `json:"document_ids"`
	UseOpenAI   bool   `json:"use_openai"`
}

func (h *PPTHandler) Generate(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Template == "" {
		req.Template = "default"
	}

	// Check if SSE stream is requested
	if c.GetHeader("Accept") == "text/event-stream" {
		h.generateStream(c, userID, req)
		return
	}

	// Regular synchronous generation
	ppt, err := h.pptService.GeneratePPT(
		c.Request.Context(),
		userID,
		req.Title,
		req.Prompt,
		req.Template,
		req.DocumentIDs,
		req.UseOpenAI,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ppt)
}

func (h *PPTHandler) generateStream(c *gin.Context, userID uint, req GenerateRequest) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// Send initial message
	fmt.Fprintf(c.Writer, "data: {\"status\":\"generating\"}\n\n")
	flusher.Flush()

	// Generate PPT (simplified for now, should use streaming AI service)
	ppt, err := h.pptService.GeneratePPT(
		c.Request.Context(),
		userID,
		req.Title,
		req.Prompt,
		req.Template,
		req.DocumentIDs,
		req.UseOpenAI,
	)

	if err != nil {
		fmt.Fprintf(c.Writer, "data: {\"status\":\"error\",\"error\":\"%s\"}\n\n", err.Error())
	} else {
		fmt.Fprintf(c.Writer, "data: {\"status\":\"completed\",\"ppt_id\":%d}\n\n", ppt.ID)
	}
	flusher.Flush()
}

func (h *PPTHandler) GetTemplates(c *gin.Context) {
	templates := h.pptService.GetTemplates()
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

type CompileRequest struct {
	LatexContent string `json:"latex_content" binding:"required"`
}

func (h *PPTHandler) Compile(c *gin.Context) {
	var req CompileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// For simplicity, create a temporary PPT record
	ppt, err := h.pptService.GeneratePPT(
		c.Request.Context(),
		userID,
		"Manual Compile",
		"Manual compilation",
		"default",
		nil,
		false,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compile the provided LaTeX
	if err := h.pptService.CompileLaTeX(ppt.ID, req.LatexContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Compilation failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, ppt)
}

func (h *PPTHandler) GetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ppts, err := h.pptService.GetPPTHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get history"})
		return
	}

	c.JSON(http.StatusOK, ppts)
}

func (h *PPTHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PPT ID"})
		return
	}

	ppt, err := h.pptService.GetPPT(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PPT not found"})
		return
	}

	c.JSON(http.StatusOK, ppt)
}

func (h *PPTHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PPT ID"})
		return
	}

	ppt, err := h.pptService.GetPPT(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PPT not found"})
		return
	}

	if ppt.PDFPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "PDF not available"})
		return
	}

	c.File(ppt.PDFPath)
}

func (h *PPTHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PPT ID"})
		return
	}

	if err := h.pptService.DeletePPT(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete PPT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PPT deleted successfully"})
}
