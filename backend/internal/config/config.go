package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Milvus   MilvusConfig
	AI       AIConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type MilvusConfig struct {
	Host string
	Port string
}

type AIConfig struct {
	OpenAIAPIKey  string
	OpenAIBaseURL string
	ClaudeAPIKey  string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type StorageConfig struct {
	UploadDir string
	OutputDir string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expireHours := 24
	if hours := os.Getenv("JWT_EXPIRE_HOURS"); hours != "" {
		if h, err := strconv.Atoi(hours); err == nil {
			expireHours = h
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "latex_ppt"),
		},
		Milvus: MilvusConfig{
			Host: getEnv("MILVUS_HOST", "localhost"),
			Port: getEnv("MILVUS_PORT", "19530"),
		},
		AI: AIConfig{
			OpenAIAPIKey:  getEnv("OPENAI_API_KEY", ""),
			OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
			ClaudeAPIKey:  getEnv("CLAUDE_API_KEY", ""),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
			ExpireHours: expireHours,
		},
		Storage: StorageConfig{
			UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
			OutputDir: getEnv("OUTPUT_DIR", "./outputs"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
