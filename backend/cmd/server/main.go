package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/api"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/config"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/internal/model"
	"github.com/qingbingwei/latex_ppt_by_claude/backend/pkg/vectordb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup database
	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.Document{},
		&model.Chunk{},
		&model.PPTRecord{},
		&model.PPTKnowledgeRef{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup Milvus collection
	milvusClient, err := vectordb.NewMilvusClient(cfg.Milvus.Host, cfg.Milvus.Port)
	if err != nil {
		log.Printf("Warning: Failed to connect to Milvus: %v", err)
	} else {
		if err := milvusClient.CreateCollection(context.Background()); err != nil {
			log.Printf("Warning: Failed to create Milvus collection: %v", err)
		}
	}

	// Ensure storage directories exist
	if err := os.MkdirAll(cfg.Storage.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := os.MkdirAll(cfg.Storage.OutputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Setup router
	router := api.SetupRouter(db, cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
