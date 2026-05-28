package config

import (
	"log"
	"os"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Env string `yaml:"env"`
	HTTP HTTP `yaml:"http_server"`
	DB DB `yaml:"db"`
	JWT JWT `yaml:"jwt"`
	Cookies Cookies `yaml:"cookies"`
	CORS CORS `yaml:"cors"`
	Auth Auth `yaml:"auth"`
}

type HTTP struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS"`
	Timeout time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	MaxBodyBytes int64 `yaml:"max_body_bytes"`
}

type DB struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	DBName string `yaml:"dbname"`
	SSLMode string `yaml:"sslmode"`
	AutoMigrate bool `yaml:"auto_migrate"`
}

type JWT struct {
	SecretKey string `yaml:"secret_key"`
	AccessTTL time.Duration `yaml:"access_ttl"`
	RefreshTTL time.Duration `yaml:"refresh_ttl"`
}

type Cookies struct {
	Secure bool `yaml:"secure"`
}

type CORS struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type Auth struct {
	RateLimitPerMinute int `yaml:"rate_limit_per_minute"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config_path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	var cfg Config
	useYAMLTags := func(c *mapstructure.DecoderConfig) {
		c.TagName = "yaml"
	}
	if err := viper.Unmarshal(&cfg, useYAMLTags); err != nil {
		log.Fatalf("error unmarshaling config: %s", err)
	}

	applyDefaults(&cfg)
	return &cfg
}

func applyDefaults(cfg *Config) {
	if cfg.HTTP.Address == "" {
		cfg.HTTP.Address = "localhost:8080"
	}
	if cfg.HTTP.Timeout == 0 {
		cfg.HTTP.Timeout = 30 * time.Second
	}
	if cfg.HTTP.IdleTimeout == 0 {
		cfg.HTTP.IdleTimeout = 60 * time.Second
	}
	if cfg.HTTP.MaxBodyBytes == 0 {
		cfg.HTTP.MaxBodyBytes = 10 << 20
	}
	if cfg.Auth.RateLimitPerMinute == 0 {
		cfg.Auth.RateLimitPerMinute = 20
	}
	if cfg.JWT.AccessTTL == 0 {
		cfg.JWT.AccessTTL = 15 * time.Minute
	}
	if cfg.JWT.RefreshTTL == 0 {
		cfg.JWT.RefreshTTL = 7 * 24 * time.Hour
	}
	if cfg.Env == "" {
		cfg.Env = "local"
	}
	if !cfg.Cookies.Secure && cfg.Env == "prod" {
		cfg.Cookies.Secure = true
	}
	if len(cfg.CORS.AllowedOrigins) == 0 {
		cfg.CORS.AllowedOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	}
}
