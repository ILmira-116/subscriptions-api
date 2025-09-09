package main

import (
	"flag"
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
	// Флаги
	runServer := flag.Bool("serve", false, "Запустить HTTP-сервер")
	migrateDB := flag.Bool("migrate", false, "Накатить миграции")
	migrateVersion := flag.String("version", "", "Версия миграции (опционально)")
	flag.Parse()

	// 1. Инициализация конфига
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// 2. Инициализация логгера
	log := logger.New(cfg)
	log.Info("Logger is ready")

	// 3. Подключение к БД
	sqlDB, err := database.InitDB(cfg)
	if err != nil {
		log.Error("Failed to connect to DB", "error", err)
		return
	}
	defer sqlDB.Close()

	// 4. Флаг миграций
	if *migrateDB {
		if *migrateVersion != "" {
			// Применяем миграции до указанной версии
			if err := database.ApplyMigrationsToVersion(sqlDB, "./migrations", *migrateVersion, log); err != nil {
				log.Error("Migration to version failed", "version", *migrateVersion, "error", err)
			}
		} else {
			// Применяем все миграции
			if err := database.ApplyAllMigrations(sqlDB, "./migrations", log); err != nil {
				log.Error("Applying all migrations failed", "error", err)
			}
		}
		return // После применения миграций сервер не запускаем
	}

	// 5. Флаг сервера
	if *runServer {
		// Репозиторий
		repo := repository.NewSubscriptionRepo(sqlDB)
		// Сервис
		svc := service.NewSubscriptionService(repo, log)
		// Роутер
		r := router.NewRouter(log, svc)
		// Сервер
		srv := server.NewServer(cfg, log, r)
		srv.Start()
		log.Info("Server is running")
		// Ожидание сигналов завершения
		shutdown.WaitForSignals(10*time.Second, log, srv)
		return
	}

	log.Info("No action specified. Use -serve to start server or -migrate to apply migrations")
}
