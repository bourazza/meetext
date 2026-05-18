package config

import "time"

type Config struct {
	App     AppConfig
	HTTP    HTTPConfig
	DB      DBConfig
	JWT     JWTConfig
	OAuth   OAuthConfig
	Storage StorageConfig
	Redis   RedisConfig
	AI      AIConfig
	Log     LogConfig
}

type AppConfig struct {
	Name        string `mapstructure:"APP_NAME"`
	Env         string `mapstructure:"APP_ENV"`
	Version     string `mapstructure:"APP_VERSION"`
	FrontendURL string `mapstructure:"FRONTEND_URL"`
}

type HTTPConfig struct {
	Host         string        `mapstructure:"HTTP_HOST"`
	Port         string        `mapstructure:"HTTP_PORT"`
	ReadTimeout  time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
}

type DBConfig struct {
	DSN          string        `mapstructure:"DATABASE_URL"`
	MaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxLifetime  time.Duration `mapstructure:"DB_MAX_LIFETIME"`
}

type JWTConfig struct {
	AccessSecret  string        `mapstructure:"JWT_ACCESS_SECRET"`
	RefreshSecret string        `mapstructure:"JWT_REFRESH_SECRET"`
	AccessTTL     time.Duration `mapstructure:"JWT_ACCESS_TTL"`
	RefreshTTL    time.Duration `mapstructure:"JWT_REFRESH_TTL"`
}

type OAuthConfig struct {
	GoogleClientID      string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret  string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL   string `mapstructure:"GOOGLE_REDIRECT_URL"`
	GitHubClientID      string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret  string `mapstructure:"GITHUB_CLIENT_SECRET"`
	GitHubRedirectURL   string `mapstructure:"GITHUB_REDIRECT_URL"`
	StateSecret         string `mapstructure:"OAUTH_STATE_SECRET"`
}

type StorageConfig struct {
	Provider  string `mapstructure:"STORAGE_PROVIDER"`
	LocalPath string `mapstructure:"STORAGE_LOCAL_PATH"`
	Bucket    string `mapstructure:"STORAGE_BUCKET"`
	Region    string `mapstructure:"STORAGE_REGION"`
	AccessKey string `mapstructure:"STORAGE_ACCESS_KEY"`
	SecretKey string `mapstructure:"STORAGE_SECRET_KEY"`
	Endpoint  string `mapstructure:"STORAGE_ENDPOINT"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"REDIS_ADDR"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type AIConfig struct {
	OllamaURL   string `mapstructure:"OLLAMA_URL"`
	OllamaModel string `mapstructure:"OLLAMA_MODEL"`
	WhisperURL  string `mapstructure:"WHISPER_URL"`
}

type LogConfig struct {
	Level     string `mapstructure:"LOG_LEVEL"`
	Pretty    bool   `mapstructure:"LOG_PRETTY"`
	File      string `mapstructure:"LOG_FILE"`
	MaxSizeMB int    `mapstructure:"LOG_MAX_SIZE_MB"`
}
