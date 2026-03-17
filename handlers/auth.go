package handlers

import (
	"auth_service/auth"
	"auth_service/config"
	"auth_service/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// Структура, хранящая все зависимости нужные для обработчиков
type AuthHandler struct {
	db            *sqlx.DB
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// Конструктор, который создает новый экземпляр AuthHandler
func NewAuthHandler(db *sqlx.DB, cfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{
		db:            db,
		accessSecret:  cfg.AccessSecret,
		refreshSecret: cfg.RefreshSecret,
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
	}
}

// Функция создания нового пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат", 400)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Все поля обязательны", 400)
		return
	}

	hash, _ := auth.HashPassword(req.Password)

	var user models.User
	err := h.db.QueryRowx(`
		INSERT INTO users (email, password_hash, created_at) 
		VALUES ($1, $2, $3) RETURNING id, email, created_at`,
		req.Email, hash, time.Now(),
	).StructScan(&user)

	if err != nil {
		http.Error(w, "Email уже используется", 409)
		return
	}

	response := models.AuthResponse{
		AccessToken:  auth.GenerateToken(user.ID, user.Email, "access", h.accessSecret, h.accessTTL),
		RefreshToken: auth.GenerateToken(user.ID, user.Email, "refresh", h.refreshSecret, h.refreshTTL),
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
}

// Функция для аутентификации существующего пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат", 400)
		return
	}

	var user models.User
	var hash string
	err := h.db.QueryRowContext(ctx, `
		SELECT id, email, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Email, &hash)

	if err != nil {
		if err == context.DeadlineExceeded {
			http.Error(w, "Сервер занят, попробуйте позже", 504)
			return
		}
		http.Error(w, "Неверный email или пароль", 401)
		return
	}

	if !auth.CheckPassword(req.Password, hash) {
		http.Error(w, "Неверный email или пароль", 401)
		return
	}

	response := models.AuthResponse{
		AccessToken:  auth.GenerateToken(user.ID, user.Email, "access", h.accessSecret, h.accessTTL),
		RefreshToken: auth.GenerateToken(user.ID, user.Email, "refresh", h.refreshSecret, h.refreshTTL),
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция принимает валидный refresh-токен и возвращает новую пару токенов
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат", 400)
		return
	}

	claims, err := auth.ValidateToken(req.RefreshToken, h.refreshSecret)
	if err != nil {
		http.Error(w, "Недействительный токен", 401)
		return
	}

	response := models.AuthResponse{
		AccessToken:  auth.GenerateToken(claims.UserID, claims.Email, "access", h.accessSecret, h.accessTTL),
		RefreshToken: auth.GenerateToken(claims.UserID, claims.Email, "refresh", h.refreshSecret, h.refreshTTL),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция-выход
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]string{"message": "Выход выполнен"})
}

// Функция возвращает данные текущего пользователя
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var user models.User
	err := h.db.QueryRowx("SELECT id, email, created_at FROM users WHERE id = $1", userID).
		StructScan(&user)

	if err != nil {
		http.Error(w, "Пользователь не найден", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
