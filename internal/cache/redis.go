package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
	keyBuilder *CacheKeyBuilder
}

// NewRedisCache creates a new Redis cache instance with optimized connection pooling
func NewRedisCache(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,

		// Connection Pool Configuration
		MaxRetries:      cfg.MaxRetries,
		PoolSize:        cfg.PoolSize,        // Maximum number of socket connections
		MinIdleConns:    cfg.MinIdleConns,    // Minimum number of idle connections
		MaxIdleConns:    cfg.MaxIdleConns,    // Maximum number of idle connections
		ConnMaxLifetime: time.Duration(cfg.ConnMaxLifetime) * time.Second, // Connection age at which client retires
		ConnMaxIdleTime: time.Duration(cfg.ConnMaxIdleTime) * time.Second, // Close idle connections after this time

		// Timeout configurations for faster failure detection
		DialTimeout:  time.Duration(cfg.DialTimeoutMs) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutMs) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeoutMs) * time.Millisecond,

		// Enable connection pooling stats for monitoring
		PoolTimeout: 4 * time.Second, // Amount of time client waits for connection if all are busy
	})

	// Test connection with shorter timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		keyBuilder: NewCacheKeyBuilder("anime-api"),
	}, nil
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	tracer := tracing.GetTracer(ctx)

	// Get connection pool stats before operation
	poolStats := r.client.PoolStats()

	ctx, span := tracer.Start(ctx, "Redis.Get",
		trace.WithAttributes(
			attribute.String("cache.operation", "get"),
			attribute.String("cache.key", key),
			attribute.String("cache.backend", "redis"),
			attribute.Int("redis.pool.hits", int(poolStats.Hits)),
			attribute.Int("redis.pool.misses", int(poolStats.Misses)),
			attribute.Int("redis.pool.timeouts", int(poolStats.Timeouts)),
			attribute.Int("redis.pool.total_conns", int(poolStats.TotalConns)),
			attribute.Int("redis.pool.idle_conns", int(poolStats.IdleConns)),
			attribute.Int("redis.pool.stale_conns", int(poolStats.StaleConns)),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		span.SetAttributes(
			attribute.Bool("cache.hit", false),
			attribute.String("cache.result", "miss"),
		)
		return nil, ErrCacheMiss
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.Bool("cache.hit", false),
			attribute.String("cache.result", "error"),
		)
		return nil, fmt.Errorf("redis get error: %w", err)
	}
	span.SetAttributes(
		attribute.Bool("cache.hit", true),
		attribute.String("cache.result", "hit"),
		attribute.Int("cache.size_bytes", len(result)),
	)
	return []byte(result), nil
}

// Set stores a value in Redis with TTL
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	tracer := tracing.GetTracer(ctx)

	// Get connection pool stats before operation
	poolStats := r.client.PoolStats()

	ctx, span := tracer.Start(ctx, "Redis.Set",
		trace.WithAttributes(
			attribute.String("cache.operation", "set"),
			attribute.String("cache.key", key),
			attribute.String("cache.backend", "redis"),
			attribute.Int("cache.ttl_seconds", int(ttl.Seconds())),
			attribute.Int("cache.size_bytes", len(value)),
			attribute.Int("redis.pool.hits", int(poolStats.Hits)),
			attribute.Int("redis.pool.misses", int(poolStats.Misses)),
			attribute.Int("redis.pool.timeouts", int(poolStats.Timeouts)),
			attribute.Int("redis.pool.total_conns", int(poolStats.TotalConns)),
			attribute.Int("redis.pool.idle_conns", int(poolStats.IdleConns)),
			attribute.Int("redis.pool.stale_conns", int(poolStats.StaleConns)),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("redis set error: %w", err)
	}
	span.SetAttributes(attribute.String("cache.result", "success"))
	return nil
}

// Delete removes a value from Redis
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "Redis.Delete",
		trace.WithAttributes(
			attribute.String("cache.operation", "delete"),
			attribute.String("cache.key", key),
			attribute.String("cache.backend", "redis"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("redis delete error: %w", err)
	}
	span.SetAttributes(attribute.String("cache.result", "success"))
	return nil
}

// DeletePattern removes all keys matching a pattern
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "Redis.DeletePattern",
		trace.WithAttributes(
			attribute.String("cache.operation", "delete_pattern"),
			attribute.String("cache.pattern", pattern),
			attribute.String("cache.backend", "redis"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("redis keys error: %w", err)
	}

	span.SetAttributes(attribute.Int("cache.keys_found", len(keys)))
	if len(keys) == 0 {
		span.SetAttributes(attribute.String("cache.result", "no_keys_found"))
		return nil
	}

	err = r.client.Del(ctx, keys...).Err()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("redis delete pattern error: %w", err)
	}
	span.SetAttributes(
		attribute.String("cache.result", "success"),
		attribute.Int("cache.keys_deleted", len(keys)),
	)
	return nil
}

// Exists checks if a key exists in Redis
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "Redis.Exists",
		trace.WithAttributes(
			attribute.String("cache.operation", "exists"),
			attribute.String("cache.key", key),
			attribute.String("cache.backend", "redis"),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, fmt.Errorf("redis exists error: %w", err)
	}
	exists := count > 0
	span.SetAttributes(
		attribute.Bool("cache.exists", exists),
		attribute.String("cache.result", "success"),
	)
	return exists, nil
}

// SetNX sets a value only if the key doesn't exist (for locking)
func (r *RedisCache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "Redis.SetNX",
		trace.WithAttributes(
			attribute.String("cache.operation", "setnx"),
			attribute.String("cache.key", key),
			attribute.String("cache.backend", "redis"),
			attribute.Int("cache.ttl_seconds", int(ttl.Seconds())),
			attribute.Int("cache.size_bytes", len(value)),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	result, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, fmt.Errorf("redis setnx error: %w", err)
	}
	span.SetAttributes(
		attribute.Bool("cache.set_success", result),
		attribute.String("cache.result", "success"),
	)
	return result, nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetKeyBuilder returns the cache key builder
func (r *RedisCache) GetKeyBuilder() *CacheKeyBuilder {
	return r.keyBuilder
}