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

	// The pool is sized to be reused rather than refilled.
	//
	// MaxIdleConns matches MaxOpenConns deliberately: Go only retains up to
	// MaxIdleConns, so anything opened above it is closed again the moment the
	// query finishes. The previous settings kept far fewer idle than they were
	// willing to open, which meant most connections were built per-query --
	// TCP, TLS and auth each time, against RDS over the internet.
	//
	// 5 is small because the database is small. db.t4g.micro allows 79
	// connections in total, and roughly 36 pods share them; a warm pool of 5
	// serves more traffic than a churning pool of 25 while claiming a fraction
	// of that budget.
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(5)

	// Connections are kept for hours because building one is expensive here.
	//
	// Measured against the production RDS instance: opening a connection costs
	// 6-52 seconds, while queries on an already-open connection cost about 7ms.
	// The gap is TLS and SCRAM authentication, both CPU-bound, on a db.t4g.micro
	// whose CPU credit balance sits at zero. Reconnecting is not a small cost
	// paid occasionally, it is the single most expensive thing the service does.
	//
	// The old 10 minute idle timeout emptied the pool during any quiet period,
	// so the next visitor after a lull paid that price -- which is exactly the
	// 45 second homepage load. Holding connections through quiet periods is
	// what removes it.
	//
	// It also preserves the statement cache. pgx v5 defaults to
	// QueryExecModeCacheStatement and caches prepared statements per
	// connection, so closing a connection discards its cached plans too. The
	// same query then re-plans from cold: 666ms on first execution against 5ms
	// on the second.
	//
	// Four hours rather than never, so a failover or DNS change is still picked
	// up without needing a restart.
	sqlDB.SetConnMaxLifetime(4 * time.Hour)
	sqlDB.SetConnMaxIdleTime(1 * time.Hour)

	// Initialize connection pool metrics collection
	poolMetrics := metrics.NewConnectionPoolMetrics(db)
	// Start collecting metrics every 30 seconds
	poolMetrics.StartPeriodicCollection(30 * time.Second)

	return &DB{DB: db}
}
