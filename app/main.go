package main

import (
	"log"
	"net/http"

	"testovoe/internal/api"
	"testovoe/internal/db"
	"testovoe/internal/handlers"
	"testovoe/internal/services"
)

func main() {
	// Инициализация базы данных
	pg, err := db.NewPostgresProvider()
	if err != nil {
		log.Fatalf("Failed to create Postgres provider: %v", err)
	}

	// Инициализация сервисов и обработчиков
	userService := services.NewUserService(pg)
	authHandler := handlers.NewAuthHandler(pg, userService)

	router := api.NewRouter(authHandler)

	defer pg.Close()

	log.Println("Сервер запущен на порту :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
