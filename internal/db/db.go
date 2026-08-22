package db

import (
	"fmt"
	"time"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/metrics"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	DB *gorm.DB
}

func NewDatabase(cfg config.DBConfig) *DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DataBase, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewTracedLogger(),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Add tracing plugin
	if err := db.Use(&TracingPlugin{}); err != nil {
		panic(fmt.Sprintf("failed to add tracing plugin: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database connection")
	}

	// Set maximum number of open connections
	// This prevents too many connections to the database
	sqlDB.SetMaxOpenConns(25)

	// Set maximum number of idle connections
	// This maintains a pool of reusable connections
	sqlDB.SetMaxIdleConns(10)

	// Set maximum lifetime of a connection
	// Kept well under any server-side idle timeout so the pool never hands out a
	// connection the server has already closed.
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Set maximum idle time for a connection
	// This helps clean up idle connections
	sqlDB.SetConnMaxIdleTime(90 * time.Second)

	// Initialize connection pool metrics collection
	poolMetrics := metrics.NewConnectionPoolMetrics(db)
	// Start collecting metrics every 30 seconds
	poolMetrics.StartPeriodicCollection(30 * time.Second)

	return &DB{DB: db}
}
