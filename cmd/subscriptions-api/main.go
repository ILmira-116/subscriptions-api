package main

import (
	"fmt"
	"time"

	"subscriptions-api/cmd/config"
	"subscriptions-api/internal/database"
	"subscriptions-api/internal/repository"
	"subscriptions-api/internal/server"
	"subscriptions-api/internal/service"
	"subscriptions-api/internal/shutdown"
	"subscriptions-api/pkg/utils/logger"
	"subscriptions-api/router"
)

func main() {
	// 1. Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Инициализация логгера
	log := logger.New(cfg)
	log.Info("Logger is ready")

	// 3. Подключаемся к БД
	sqlDB, err := database.InitDB(cfg)
	if err != nil {
		log.Error("Failed to connect to DB", "error", err)
		return
	}
	defer sqlDB.Close()

	// 4. Применяем миграции
	if err := database.ApplyMigrations(sqlDB, "./migrations", log); err != nil {
		log.Error("Migration failed", "error", err)
		return
	}

	// 5. Репозиторий
	repo := repository.NewSubscriptionRepo(sqlDB)

	// 6. Сервис
	svc := service.NewSubscriptionService(repo, log)

	// 7. Инициализация и запуск роутера
	r := router.NewRouter(log, svc)

	// 8. Создаём сервер и запускаем его
	srv := server.NewServer(cfg, log, r)
	srv.Start()

	log.Info("Server is running")

	// 9. Shutdown при сигнале
	shutdown.WaitForSignals(10*time.Second, log, srv)

}
