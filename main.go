package main

import (
	"auth_service/config"
	"auth_service/database"
	"auth_service/handlers"
	"auth_service/middleware"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	defer db.Close()

	database.Migrate(db)

	authHandler := handlers.NewAuthHandler(db, cfg.JWT)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/auth/logout", authHandler.Logout)

	protected := middleware.Auth(cfg.JWT.AccessSecret)
	mux.Handle("GET /api/auth/me", protected(http.HandlerFunc(authHandler.Me)))

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: mux,
	}

	go func() {
		log.Printf("Сервер запущен на :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println("Сервер остановлен")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
