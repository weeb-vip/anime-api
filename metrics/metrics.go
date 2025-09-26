package metrics

import (
	"github.com/weeb-vip/anime-api/config"
	metricsLib "github.com/weeb-vip/go-metrics-lib"
	"github.com/weeb-vip/go-metrics-lib/clients/prometheus"
)

var metricsInstance metricsLib.MetricsImpl

var prometheusInstance *prometheus.PrometheusClient

func NewMetricsInstance() metricsLib.MetricsImpl {
	if metricsInstance == nil {
		prometheusInstance = NewPrometheusInstance()
		initMetrics(prometheusInstance)
		metricsInstance = metricsLib.NewMetrics(prometheusInstance, 1)
	}
	return metricsInstance
}

func NewPrometheusInstance() *prometheus.PrometheusClient {
	if prometheusInstance == nil {
		prometheusInstance = prometheus.NewPrometheusClient()
		initMetrics(prometheusInstance)
	}
	return prometheusInstance
}

func initMetrics(prometheusInstance *prometheus.PrometheusClient) {
	prometheusInstance.CreateHistogramVec("resolver_request_duration_histogram_milliseconds", "graphql resolver millisecond", []string{"service", "protocol", "resolver", "result", "env"}, []float64{
		// create buckets 10000 split into 10 buckets
		100,
		200,
		300,
		400,
		500,
		600,
		700,
		800,
		900,
		1000,
	})

	prometheusInstance.CreateHistogramVec("database_query_duration_histogram_milliseconds", "database calls millisecond", []string{"service", "table", "method", "result", "env"}, []float64{
		// create buckets 10000 split into 10 buckets
		100,
		200,
		300,
		400,
		500,
		600,
		700,
		800,
		900,
		1000,
	})

	// Database connection pool metrics
	prometheusInstance.CreateGaugeVec("database_connection_pool_open_connections", "Number of open database connections", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("database_connection_pool_in_use_connections", "Number of database connections in use", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("database_connection_pool_idle_connections", "Number of idle database connections", []string{"service", "env"})
	prometheusInstance.CreateCounterVec("database_connection_pool_wait_total", "Total number of connection waits", []string{"service", "env"})
	prometheusInstance.CreateHistogramVec("database_connection_acquisition_duration_milliseconds", "Time to acquire database connection", []string{"service", "env"}, []float64{
		1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000,
	})

	// Redis connection pool metrics
	prometheusInstance.CreateGaugeVec("redis_connection_pool_hits_total", "Number of times connection pool hit", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("redis_connection_pool_misses_total", "Number of times connection pool missed", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("redis_connection_pool_timeouts_total", "Number of connection pool timeouts", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("redis_connection_pool_total_connections", "Total number of connections in pool", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("redis_connection_pool_idle_connections", "Number of idle connections in pool", []string{"service", "env"})
	prometheusInstance.CreateGaugeVec("redis_connection_pool_stale_connections", "Number of stale connections in pool", []string{"service", "env"})
}

func GetCurrentEnv() string {
	cfg := config.LoadConfigOrPanic()
	return cfg.AppConfig.Env
}
