package main

import (
	"github.com/RVodassa/geo-microservices-proxy/app"
	myLogger "github.com/RVodassa/slog_utils/slog_logger"
	"log"
	"log/slog"
	"os"
)

// @title Geo Microservices API
// @version 1.0
// @description API для работы с геоданными

// @contact.name API Support
// @contact.email support@geo.com

// @license.name Apache 2.0

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
// @host localhost:8080

// TODO: Доделать ветки логирования и возврата ошибок
// TODO: Передача через интерфейсы
// TODO: Тесты

const (
	defaultConfigPath = "/app/proxy-service/configs/config.yaml" // Путь в контейнере
	configEnvVar      = "CONFIG_PATH"                            // EXPORT CONFIG_PATH
)

func main() {
	const op = "main.main"

	configPath := os.Getenv(configEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
		log.Printf("Путь к конфигурации по умолчанию: %s", configPath)
	}

	cfg, err := app.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	logger := setupLogger(cfg)
	log.Printf("Логер инициализирван. LogLevel=%s; Env=%s", cfg.LogLevel, cfg.Env)
	// Инициализация и запуск приложения
	newApp := app.NewApp(logger, cfg)
	if err = newApp.Run(); err != nil {
		logger.Error(op, err)
		os.Exit(1)
	}
}

func setupLogger(cfg *app.Config) *slog.Logger {
	var logLevel slog.Level

	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	logger := myLogger.SetupLogger(cfg.Env, logLevel)
	return logger
}
