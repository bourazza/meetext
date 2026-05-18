package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("env")
	v.AutomaticEnv()

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if !strings.Contains(err.Error(), "no such file") {
				return nil, fmt.Errorf("config: read error: %w", err)
			}
		}
	}

	cfg := &Config{}

	cfg.App = AppConfig{
		Name:        v.GetString("APP_NAME"),
		Env:         v.GetString("APP_ENV"),
		Version:     v.GetString("APP_VERSION"),
		FrontendURL: v.GetString("FRONTEND_URL"),
	}

	cfg.HTTP = HTTPConfig{
		Host:         v.GetString("HTTP_HOST"),
		Port:         v.GetString("HTTP_PORT"),
		ReadTimeout:  v.GetDuration("HTTP_READ_TIMEOUT"),
		WriteTimeout: v.GetDuration("HTTP_WRITE_TIMEOUT"),
		IdleTimeout:  v.GetDuration("HTTP_IDLE_TIMEOUT"),
	}

	cfg.DB = DBConfig{
		DSN:          v.GetString("DATABASE_URL"),
		MaxOpenConns: v.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns: v.GetInt("DB_MAX_IDLE_CONNS"),
		MaxLifetime:  v.GetDuration("DB_MAX_LIFETIME"),
	}

	cfg.JWT = JWTConfig{
		AccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
		RefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
		AccessTTL:     v.GetDuration("JWT_ACCESS_TTL"),
		RefreshTTL:    v.GetDuration("JWT_REFRESH_TTL"),
	}

	cfg.OAuth = OAuthConfig{
		GoogleClientID:     v.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: v.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  v.GetString("GOOGLE_REDIRECT_URL"),
		GitHubClientID:     v.GetString("GITHUB_CLIENT_ID"),
		GitHubClientSecret: v.GetString("GITHUB_CLIENT_SECRET"),
		GitHubRedirectURL:  v.GetString("GITHUB_REDIRECT_URL"),
		StateSecret:        v.GetString("OAUTH_STATE_SECRET"),
	}

	cfg.Storage = StorageConfig{
		Provider:  v.GetString("STORAGE_PROVIDER"),
		LocalPath: v.GetString("STORAGE_LOCAL_PATH"),
		Bucket:    v.GetString("STORAGE_BUCKET"),
		Region:    v.GetString("STORAGE_REGION"),
		AccessKey: v.GetString("STORAGE_ACCESS_KEY"),
		SecretKey: v.GetString("STORAGE_SECRET_KEY"),
		Endpoint:  v.GetString("STORAGE_ENDPOINT"),
	}

	cfg.Redis = RedisConfig{
		Addr:     v.GetString("REDIS_ADDR"),
		Password: v.GetString("REDIS_PASSWORD"),
		DB:       v.GetInt("REDIS_DB"),
	}

	cfg.AI = AIConfig{
		OllamaURL:   v.GetString("OLLAMA_URL"),
		OllamaModel: v.GetString("OLLAMA_MODEL"),
		WhisperURL:  v.GetString("WHISPER_URL"),
	}

	cfg.Log = LogConfig{
		Level:     v.GetString("LOG_LEVEL"),
		Pretty:    v.GetBool("LOG_PRETTY"),
		File:      v.GetString("LOG_FILE"),
		MaxSizeMB: v.GetInt("LOG_MAX_SIZE_MB"),
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "meetext")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_VERSION", "0.1.0")
	v.SetDefault("FRONTEND_URL", "http://localhost:3000")

	v.SetDefault("HTTP_HOST", "0.0.0.0")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("HTTP_READ_TIMEOUT", 15*time.Second)
	v.SetDefault("HTTP_WRITE_TIMEOUT", 15*time.Second)
	v.SetDefault("HTTP_IDLE_TIMEOUT", 60*time.Second)

	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_MAX_IDLE_CONNS", 5)
	v.SetDefault("DB_MAX_LIFETIME", 5*time.Minute)

	v.SetDefault("JWT_ACCESS_TTL", 15*time.Minute)
	v.SetDefault("JWT_REFRESH_TTL", 7*24*time.Hour)

	v.SetDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/google/callback")
	v.SetDefault("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth/github/callback")
	v.SetDefault("OAUTH_STATE_SECRET", "change-me-oauth-state-secret")

	v.SetDefault("STORAGE_PROVIDER", "local")
	v.SetDefault("STORAGE_LOCAL_PATH", "./uploads")

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("OLLAMA_URL", "http://localhost:11434")
	v.SetDefault("OLLAMA_MODEL", "llama3")
	v.SetDefault("WHISPER_URL", "http://localhost:9000")

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_PRETTY", true)
	v.SetDefault("LOG_FILE", "./logs/meetext.log")
	v.SetDefault("LOG_MAX_SIZE_MB", 100)
}
