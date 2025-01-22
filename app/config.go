package app

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	LogLevel     string        `yaml:"logLevel" env-default:"info"`
	GRPCServices []GRPCService `yaml:"grpc_services"`
	HTTPPort     string        `yaml:"http_port"`
}

type GRPCService struct {
	Name string `yaml:"name"`
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %v", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %v", err)
	}

	log.Printf("Загружена конфигурация: %+v", cfg)
	return &cfg, nil
}
