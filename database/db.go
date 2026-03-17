package database

import (
	"auth_service/config"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Функция для подключения к БД
func NewDB(cfg config.DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("БД не отвечает: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	var one int
	err = db.QueryRow("SELECT 1").Scan(&one)
	if err != nil {
		return nil, fmt.Errorf("БД не отвечает на запросы: %w", err)
	}

	log.Println("БД готова к работе")
	return db, nil
}

// Функция создания таблицы, при их отсутсвии
func Migrate(db *sqlx.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);
	
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		token_hash TEXT UNIQUE NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		revoked BOOLEAN DEFAULT FALSE
	);`

	db.MustExec(schema)
	log.Println("Миграции применены")
}
