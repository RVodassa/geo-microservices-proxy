package main

import (
	"github.com/RVodassa/geo-microservices-proxy/app"
	"log"
	"os"
)

const (
	defaultConfigPath = "/app/configs/config.yaml" // Путь в контейнере
	configEnvVar      = "CONFIG_PATH"
)

func main() {
	// Путь к конфигурации
	configPath := os.Getenv(configEnvVar)
	if configPath == "" {
		configPath = defaultConfigPath
		log.Printf("Используется конфигурация по умолчанию: %s", configPath)
	}

	// Запуск приложения
	if err := app.RunApp(configPath); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
