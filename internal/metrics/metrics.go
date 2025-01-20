package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Метрики prometheus
var (
	// EndpointRequestCount метрика-счетчик обращений к эндпоинтам
	EndpointRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geo_service_requests_total",
			Help: "Количество запросов к каждому эндпоинту",
		},
		[]string{"endpoint"},
	)
	// EndpointRequestDuration метрика времени на запросы к эндпоинтам
	EndpointRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "geo_service_request_duration_seconds",
			Help:    "Время выполнения запросов к эндпоинтам",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
	// CacheDuration метрика времени ответа от cache
	CacheDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cache_access_duration_seconds",
			Help: "время ответа редис",
		},
		[]string{"method"},
	)
	// ApiDuration метрика времени ответа API
	ApiDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "api_access_duration_seconds",
			Help: "время ответа api",
		},
		[]string{"method"},
	)
)
