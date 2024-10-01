package app

import (
	"github.com/VadimBorzenkov/online-song-library/config"
	"github.com/VadimBorzenkov/online-song-library/internal/db"
	"github.com/VadimBorzenkov/online-song-library/internal/delivery/handler"
	"github.com/VadimBorzenkov/online-song-library/internal/delivery/routes"
	"github.com/VadimBorzenkov/online-song-library/internal/log"
	"github.com/VadimBorzenkov/online-song-library/internal/repository"
	"github.com/VadimBorzenkov/online-song-library/internal/service"
	"github.com/VadimBorzenkov/online-song-library/pkg/migrator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Run() {
	// Инициализация логгера
	logger := log.InitLogger()

	// Загрузка конфигурации
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация базы данных
	dbase := db.Init(config)
	defer func() {
		if err := db.Close(dbase); err != nil {
			logger.Errorf("Failed to close database: %v", err)
		}
	}()

	// Выполнение миграций
	if err := migrator.RunDatabaseMigrations(dbase); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация репозитория
	repo := repository.NewApiRepository(dbase, logger)

	// Инициализация сервисного уровня
	svc := service.NewApiSetvice(repo, logger, config)

	// Инициализация хендлеров
	handler := handler.NewApiHandler(svc, logger)

	// Инициализация приложения Fiber
	app := fiber.New()

	// Настройка CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// Регистрация маршрутов
	routes.RegistrationRoutes(app, handler)

	// Запуск сервера
	logger.Infof("Starting server on port %s", config.Port)
	if err := app.Listen(config.Port); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
