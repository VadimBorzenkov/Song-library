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
)

func Run() {
	logger := log.InitLogger()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	dbase := db.Init(config)
	defer func() {
		if err := db.Close(dbase); err != nil {
			logger.Errorf("Failed to close database: %v", err)
		}
	}()

	if err := migrator.RunDatabaseMigrations(dbase); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewApiRepository(dbase, logger)

	svc := service.NewApiService(repo, logger, config)

	handler := handler.NewApiHandler(svc, logger)

	app := fiber.New()

	routes.RegistrationRoutes(app, handler)

	logger.Infof("Starting server on port %s", config.Port)
	if err := app.Listen(config.Port); err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
}
