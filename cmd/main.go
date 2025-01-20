package main

import (
	"errors"
	httpController "github.com/RVodassa/geo-microservices-proxy/internal/controller"
	"github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/cache"
	grpc_auth_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_auth-service"
	grpc_geo_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_geo-service"
	grpc_user_service "github.com/RVodassa/geo-microservices-proxy/internal/infrastructure/grpc_user-service"
	_ "github.com/RVodassa/geo-microservices-proxy/pkg/swagger"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/RVodassa/geo-microservices-proxy/internal/service"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func main() {
	// TODO: Вынести все настройки в конфиг
	// TODO: logs
	// TODO: настроить метрики
	// TODO: тесты

	connGeoService, err := setupGRPCClient("geo-service:30303")
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer connGeoService.Close()

	connAuthService, err := setupGRPCClient("auth-service:20202")
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer connAuthService.Close()

	connUserService, err := setupGRPCClient("user-service:10101")
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer connUserService.Close()

	cacheService := cache.NewCacheService()
	authService := grpc_auth_service.NewAuthServiceClient(connAuthService)
	geoService := grpc_geo_service.NewGeoServiceClient(connGeoService)
	userService := grpc_user_service.NewUserServiceClient(connUserService)
	proxyService := service.NewProxyGeoService(cacheService, authService, geoService, userService)

	controller := httpController.NewHttpController(proxyService)

	r := chi.NewRouter()

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", controller.RegisterHandler)
		r.Post("/login", controller.LoginHandler)
	})

	r.With(controller.JwtMiddleware).Route("/api/user", func(r chi.Router) {
		r.Post("/list", controller.ListUsersHandler)
		r.Post("/profile", controller.ProfileHandler)
	})

	r.With(controller.JwtMiddleware).Route("/api/address", func(r chi.Router) {
		r.Post("/geocode", controller.GeoCodeHandler)
		r.Post("/search", controller.SearchHandler)
	})

	// Маршрут для Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // URL для swagger.json
	))

	// CTRL + C = graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Конфигурация HTTP-сервера
	port := ":8080"
	httpSrv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	// Запуск сервера в горутине
	errChan := make(chan error, 1) // канал для получения ошибки
	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		if err := httpSrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				errChan <- err
				return
			}
		}
	}()
	defer httpSrv.Close()

	// Ловим сигналы об ошибках и остановке сервера
	select {
	case err = <-errChan:
		log.Println("Ошибка запуска сервера:", err)
		return
	}
}

func setupGRPCClient(address string) (*grpc.ClientConn, error) {
	return grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
