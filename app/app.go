package app

import (
	"fmt"
	httpController "github.com/RVodassa/geo-microservices-proxy/internal/controller"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/logger"
	"github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/cache"
	"github.com/RVodassa/geo-microservices-proxy/internal/serve"
	"github.com/RVodassa/geo-microservices-proxy/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log logger.Logger
	cfg *Config
}

func NewApp(log logger.Logger, cfg *Config) *App {
	return &App{
		log: log,
		cfg: cfg,
	}
}

func (app *App) Run() error {
	const op = "app.Run"

	// TODO: тесты
	// TODO: настроить метрики

	app.log.Info("установка соединений с gRPC сервисами...", "op", op)
	services, err := app.setupGRPCClient(app.cfg.GRPCServices)
	if err != nil {
		return fmt.Errorf("%s: ошибка подключения к gRPC сервисам: %v", op, err)
	}

	defer func() {
		app.log.Info("закрытие gRPC соединений...", "op", op)
		for name, conn := range services {
			if err = conn.Close(); err != nil {
				app.log.Error("ошибка при закрытии соединения", "op", op, "service", name, "error", err)
			}
		}
	}()

	app.log.Info("инициализация сервисов...", "op", op)
	cacheService := cache.NewCacheService()
	if cacheService == nil {
		return fmt.Errorf("%s: ошибка инициализации cacheService", op)
	}

	proxyService := service.NewProxyGeoService(app.log, cacheService, services)
	if proxyService == nil {
		return fmt.Errorf("%s: ошибка инициализации proxyService", op)
	}

	controller := httpController.NewHttpController(app.log, proxyService)
	if controller == nil {
		return fmt.Errorf("%s: ошибка инициализации controller", op)
	}

	port := app.cfg.HTTPPort
	if port == "" {
		return fmt.Errorf("%s: отсутствует port", op)
	}

	r := httpController.NewRouter(controller)

	newServe := serve.NewServe(port, r, app.log)
	if err = newServe.RunServe(); err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	app.log.Info("HTTP-сервер запущен", "op", op, "port", port)
	return nil
}

func (app *App) setupGRPCClient(grpcServices []GRPCService) (map[string]*grpc.ClientConn, error) {
	const op = "app.setupGRPCClient"

	connections := make(map[string]*grpc.ClientConn) // service_name:client

	for _, grpcServ := range grpcServices {
		addr := fmt.Sprintf("%s:%d", grpcServ.Addr, grpcServ.Port)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			for _, c := range connections {
				_ = c.Close()
			}
			return nil, fmt.Errorf("%s: ошибка подключения к %s: %v", op, addr, err)
		}

		if _, exists := connections[grpcServ.Name]; exists {
			return nil, fmt.Errorf("%s: дубликат имени сервиса: %s", op, grpcServ.Name)
		}

		connections[grpcServ.Name] = conn
	}

	app.log.Info("Успешное подключение к gRPC сервисам", "op", op)
	return connections, nil
}
