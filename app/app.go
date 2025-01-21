package app

import (
	"fmt"
	httpController "github.com/RVodassa/geo-microservices-proxy/internal/controller"
	"github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/cache"
	"github.com/RVodassa/geo-microservices-proxy/internal/serve"
	"github.com/RVodassa/geo-microservices-proxy/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

func RunApp(configPath string) error {
	// TODO: logs
	// TODO: тесты
	// TODO: настроить метрики

	log.Println("Загрузка конфигурации...")
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %v", err)
	}
	log.Printf("Конфигурация загружена: %+v", cfg)

	log.Println("Установка соединений с gRPC сервисами...")
	grpcConns, err := setupGRPCClient(cfg.GRPCServices)
	if err != nil {
		return fmt.Errorf("ошибка подключения к gRPC сервисам: %v", err)
	}
	defer func() {
		log.Println("Закрытие gRPC соединений...")
		for name, conn := range grpcConns {
			if err = conn.Close(); err != nil {
				log.Printf("Ошибка при закрытии соединения с %s: %v", name, err)
			}
		}
	}()

	log.Println("Инициализация сервисов...")
	cacheService := cache.NewCacheService()
	if cacheService == nil {
		return fmt.Errorf("ошибка инициализации cacheService")
	}

	proxyService := service.NewProxyGeoService(cacheService, grpcConns)
	if proxyService == nil {
		return fmt.Errorf("ошибка инициализации proxyService")
	}

	controller := httpController.NewHttpController(proxyService)
	if controller == nil {
		return fmt.Errorf("ошибка инициализации controller")
	}

	log.Println("Запуск HTTP-сервера...")
	r := httpController.NewRouter(controller)
	if err = serve.RunServe(cfg.HTTPPort, r); err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}
	return nil
}

func setupGRPCClient(grpcServices []GRPCService) (map[string]*grpc.ClientConn, error) {

	conns := make(map[string]*grpc.ClientConn) // service_name:client

	for _, grpcServ := range grpcServices {
		addr := fmt.Sprintf("%s:%d", grpcServ.Addr, grpcServ.Port)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			for _, c := range conns {
				_ = c.Close()
			}
			return nil, fmt.Errorf("ошибка подключения к %s: %v", addr, err)
		}
		conns[grpcServ.Name] = conn
		log.Printf("Успешно подключен к gRPC сервису: %s", addr)
	}

	return conns, nil
}
