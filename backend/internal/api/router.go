package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/api/handler"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/api/middleware"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/config"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/repository"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/service"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/ai"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/embedding"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/latex"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/vectordb"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set mode
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Initialize clients
	var openaiClient *ai.OpenAIClient
	if cfg.AI.OpenAIAPIKey != "" {
		openaiClient = ai.NewOpenAIClient(cfg.AI.OpenAIAPIKey, cfg.AI.OpenAIBaseURL)
	}

	var claudeClient *ai.ClaudeClient
	if cfg.AI.ClaudeAPIKey != "" {
		claudeClient = ai.NewClaudeClient(cfg.AI.ClaudeAPIKey)
	}

	embeddingClient := embedding.NewOpenAIEmbedding(cfg.AI.OpenAIAPIKey, cfg.AI.OpenAIBaseURL)

	milvusClient, err := vectordb.NewMilvusClient(cfg.Milvus.Host, cfg.Milvus.Port)
	if err != nil {
		panic("Failed to connect to Milvus: " + err.Error())
	}

	latexCompiler := latex.NewCompiler(cfg.Storage.OutputDir)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	docRepo := repository.NewDocumentRepository(db)
	pptRepo := repository.NewPPTRepository(db)

	// Initialize services
	knowledgeService := service.NewKnowledgeService(docRepo, embeddingClient, milvusClient, cfg.Storage.UploadDir)
	aiService := service.NewAIService(openaiClient, claudeClient)
	pptService := service.NewPPTService(pptRepo, knowledgeService, aiService, latexCompiler, cfg.Storage.OutputDir)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeService)
	pptHandler := handler.NewPPTHandler(pptService)

	// Public routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Check)

		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.Auth(cfg.JWT.Secret))
	{
		// Auth
		protected.GET("/auth/profile", authHandler.GetProfile)

		// Knowledge base
		knowledge := protected.Group("/knowledge")
		{
			knowledge.POST("/upload", knowledgeHandler.Upload)
			knowledge.GET("/list", knowledgeHandler.List)
			knowledge.GET("/:id", knowledgeHandler.Get)
			knowledge.DELETE("/:id", knowledgeHandler.Delete)
			knowledge.POST("/search", knowledgeHandler.Search)
		}

		// PPT generation
		ppt := protected.Group("/ppt")
		{
			ppt.POST("/generate", pptHandler.Generate)
			ppt.GET("/templates", pptHandler.GetTemplates)
			ppt.POST("/compile", pptHandler.Compile)
			ppt.GET("/history", pptHandler.GetHistory)
			ppt.GET("/:id", pptHandler.Get)
			ppt.GET("/:id/download", pptHandler.Download)
			ppt.DELETE("/:id", pptHandler.Delete)
		}
	}

	return router
}
