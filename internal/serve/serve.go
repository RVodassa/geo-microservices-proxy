package serve

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// RunServe запускает HTTP-сервер с поддержкой graceful shutdown.
// port - порт, на котором будет запущен сервер (например, ":8080").
// handler - обработчик запросов (например, *chi.Mux).
// Возвращает ошибку, если сервер не удалось запустить или завершить.
func RunServe(port string, handler http.Handler) error {
	// Настройка HTTP-сервера
	httpSrv := &http.Server{
		Addr:    port,
		Handler: handler,
		// Дополнительные настройки (опционально)
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Запуск сервера в горутине
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание graceful shutdown
	if err := waitForShutdown(httpSrv); err != nil {
		return fmt.Errorf("ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Сервер успешно завершил работу")
	return nil
}

// waitForShutdown ожидает сигналов завершения и выполняет graceful shutdown сервера.
// Возвращает ошибку, если shutdown не удался.
func waitForShutdown(server *http.Server) error {
	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала завершения
	<-shutdown
	log.Println("Завершение работы сервера...")

	// Graceful shutdown для HTTP-сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка при завершении работы сервера: %v", err)
	}

	return nil
}
