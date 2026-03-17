package models

import (
	"time"
)

// Структура для хранения refresh токена в БД
type RefreshToken struct {
	ID int `db:"id" json:"-"`

	UserID int `db:"user_id" json:"-"`

	TokenHash string `db:"token_hash" json:"-"`

	ExpiresAt time.Time `db:"expires_at" json:"-"`

	Revoked   bool      `db:"revoked" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// Структура для ответа клиенту
type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Структура которая создает данные для запроса на создание refresh токена
type CreateRefreshTokenRequest struct {
	UserID int `json:"user_id"`
}

// Структура которая создает данные для запроса на проверку refresh токена
type ValidateRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Структура информация о refresh токене
type RefreshTokenInfo struct {
	UserID    int
	TokenID   string
	ExpiresAt time.Time
	IsValid   bool
}
