package metrics

import (
	"time"

	metricsLib "github.com/weeb-vip/go-metrics-lib"
	"gorm.io/gorm"
)

type ConnectionPoolMetrics struct {
	metrics metricsLib.MetricsImpl
	db      *gorm.DB
	env     string
}

func NewConnectionPoolMetrics(db *gorm.DB) *ConnectionPoolMetrics {
	return &ConnectionPoolMetrics{
		metrics: NewMetricsInstance(),
		db:      db,
		env:     GetCurrentEnv(),
	}
}

// UpdateMetrics updates connection pool metrics
func (cpm *ConnectionPoolMetrics) UpdateMetrics() error {
	sqlDB, err := cpm.db.DB()
	if err != nil {
		return err
	}

	stats := sqlDB.Stats()

	// Update gauge metrics
	cpm.metrics.GaugeMetric("database_connection_pool_open_connections", float64(stats.OpenConnections), map[string]string{
		"service": "anime-api",
		"env":     cpm.env,
	})

	cpm.metrics.GaugeMetric("database_connection_pool_in_use_connections", float64(stats.InUse), map[string]string{
		"service": "anime-api",
		"env":     cpm.env,
	})

	cpm.metrics.GaugeMetric("database_connection_pool_idle_connections", float64(stats.Idle), map[string]string{
		"service": "anime-api",
		"env":     cpm.env,
	})

	return nil
}

// RecordConnectionAcquisition records the time it took to acquire a connection
func (cpm *ConnectionPoolMetrics) RecordConnectionAcquisition(duration time.Duration) {
	cpm.metrics.HistogramMetric("database_connection_acquisition_duration_milliseconds", float64(duration.Milliseconds()), map[string]string{
		"service": "anime-api",
		"env":     cpm.env,
	})
}

// RecordConnectionWait records when a connection wait occurs
func (cpm *ConnectionPoolMetrics) RecordConnectionWait() {
	cpm.metrics.CountMetric("database_connection_pool_wait_total", map[string]string{
		"service": "anime-api",
		"env":     cpm.env,
	})
}

// StartPeriodicCollection starts periodic collection of connection pool metrics
func (cpm *ConnectionPoolMetrics) StartPeriodicCollection(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			cpm.UpdateMetrics()
		}
	}()
}