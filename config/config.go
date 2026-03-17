package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Главная структура проекта
type Config struct {
	DB     DBConfig
	JWT    JWTConfig
	Server ServerConfig
}

// Структура для базы данных
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Структура для настройки JWT-токенов
type JWTConfig struct {
	AccessSecret  string // ← добавить
	RefreshSecret string // ← добавить
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

// Структура настройки HTTP-сервера
type ServerConfig struct {
	Port string
}

// Функция наполняет главную структуру Config
func Load() *Config {
	godotenv.Load()

	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "auth"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret"),
			AccessTTL:     parseDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL:    parseDuration("JWT_REFRESH_TTL", 24*7*time.Hour),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
	}
}

// Функция для получения переменных окружения
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Фунцкция, которая достает настройки времени из переменных окружения
func parseDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
