package cache

import (
	"context"
	"encoding/json"
	"github.com/RVodassa/geo-microservices-proxy/internal/domain/entity"
	"github.com/RVodassa/geo-microservices-proxy/internal/metrics"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var MetricsCacheDuration = metrics.CacheDuration

// cacheMetrics хранит метрики для кеша
type cacheMetrics struct {
}

// recordDuration записывает время выполнения функции в метрику
func (m *cacheMetrics) recordDuration(method string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start).Seconds()
	MetricsCacheDuration.WithLabelValues(method).Observe(duration)
	return err
}

// recordDurationWithResult записывает время выполнения функции с возвращаемым значением в метрику
func (m *cacheMetrics) recordDurationWithResult(method string, fn func() ([]byte, error)) ([]byte, error) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start).Seconds()
	MetricsCacheDuration.WithLabelValues(method).Observe(duration)
	return result, err
}

// CacheServiceProvider интерфейс для работы с кэшированием
type CacheServiceProvider interface {
	Set(ctx context.Context, key string, addresses []*entity.Address) error
	Get(ctx context.Context, key string) (entity.ResponseBody, error)
}

type cacheService struct {
	client  *redis.Client
	metrics *cacheMetrics
}

// newRedisClient создает новый клиент Redis
func newRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	return client
}

// NewCacheService создает новый экземпляр кэш-сервиса
func NewCacheService() CacheServiceProvider {
	return &cacheService{
		client: newRedisClient(),
	}
}

// Set метод для сохранения данных в Redis
func (c *cacheService) Set(ctx context.Context, key string, addresses []*entity.Address) error {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return err
	}

	addressJson, err := json.Marshal(addresses)
	if err != nil {
		return err
	}

	return c.metrics.recordDuration("SET", func() error {
		if err := c.client.Set(ctx, key, addressJson, 0).Err(); err != nil {
			log.Printf("Ошибка при записи в Redis: %v\n", err)
			return err
		}
		log.Println("Данные успешно сохранены в Redis")
		return nil
	})
}

// Get метод для получения данных из Redis
func (c *cacheService) Get(ctx context.Context, key string) (entity.ResponseBody, error) {
	if err := c.client.Ping(ctx).Err(); err != nil {
		return entity.ResponseBody{}, err
	}

	var addresses []*entity.Address

	data, err := c.metrics.recordDurationWithResult("GET", func() ([]byte, error) {
		return c.client.Get(ctx, key).Bytes()
	})
	if err != nil {
		return entity.ResponseBody{}, err
	}

	if err := json.Unmarshal(data, &addresses); err != nil {
		log.Printf("Ошибка десериализации данных: %v\n", err)
		return entity.ResponseBody{}, err
	}

	resp := entity.ResponseBody{
		Addresses: addresses,
	}

	return resp, nil
}
