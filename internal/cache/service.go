package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/metrics"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"github.com/goccy/go-json"
)

// CacheService provides high-level caching operations with JSON serialization
type CacheService struct {
	cache      Cache
	keyBuilder *CacheKeyBuilder
	config     config.RedisConfig
}

// NewCacheService creates a new cache service with high-performance JSON
func NewCacheService(cache Cache, cfg config.RedisConfig) *CacheService {
	return &CacheService{
		cache:      cache,
		keyBuilder: GetKeyBuilder(),
		config:     cfg,
	}
}

// GetJSON retrieves and unmarshals JSON data from cache
func (c *CacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "CacheService.GetJSON",
		trace.WithAttributes(
			attribute.String("cache.operation", "get_json"),
			attribute.String("cache.key", key),
			attribute.String("cache.layer", "service"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	// Phase 1: Redis Get Operation
	redisStartTime := time.Now()
	data, err := c.cache.Get(ctx, key)
	redisEndTime := time.Now()
	redisDuration := redisEndTime.Sub(redisStartTime)

	span.SetAttributes(
		attribute.Int64("cache.redis_duration_us", redisDuration.Microseconds()),
		attribute.Int("cache.data_size_bytes", len(data)),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err == ErrCacheMiss {
			span.SetAttributes(attribute.String("cache.result", "miss"))
			metrics.GetAppMetrics().DatabaseMetric(
				float64(time.Since(startTime).Milliseconds()),
				"cache",
				"get",
				"miss",
			)
		} else {
			span.SetAttributes(attribute.String("cache.result", "error"))
			metrics.GetAppMetrics().DatabaseMetric(
				float64(time.Since(startTime).Milliseconds()),
				"cache",
				"get",
				metrics.Error,
			)
		}
		return err
	}

	// Phase 2: JSON Unmarshaling (using goccy/go-json for maximum performance)
	unmarshalStartTime := time.Now()
	err = json.Unmarshal(data, dest)
	unmarshalEndTime := time.Now()
	unmarshalDuration := unmarshalEndTime.Sub(unmarshalStartTime)

	span.SetAttributes(
		attribute.Int64("cache.unmarshal_duration_us", unmarshalDuration.Microseconds()),
		attribute.Int64("cache.total_duration_us", time.Since(startTime).Microseconds()),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("unmarshal error: %v", err))
		span.SetAttributes(attribute.String("cache.result", "unmarshal_error"))
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"cache",
			"get",
			metrics.Error,
		)
		return fmt.Errorf("cache unmarshal error: %w", err)
	}

	span.SetAttributes(
		attribute.String("cache.result", "hit"),
		attribute.Bool("cache.success", true),
	)
	span.SetStatus(codes.Ok, "cache hit")

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"get",
		"hit",
	)
	return nil
}

// SetJSON marshals and stores JSON data in cache
func (c *CacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	// Run cache set asynchronously - fire and forget
	go func() {
		startTime := time.Now()
		err := c.cache.Set(context.Background(), key, data, ttl)

		result := metrics.Success
		if err != nil {
			result = metrics.Error
		}

		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"cache",
			"set",
			result,
		)
	}()

	return nil
}

// Delete removes a key from cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	startTime := time.Now()

	err := c.cache.Delete(ctx, key)

	result := metrics.Success
	if err != nil {
		result = metrics.Error
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"delete",
		result,
	)
	return err
}

// DeletePattern removes all keys matching a pattern
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	startTime := time.Now()

	err := c.cache.DeletePattern(ctx, pattern)

	result := metrics.Success
	if err != nil {
		result = metrics.Error
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"delete_pattern",
		result,
	)
	return err
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	return c.cache.Exists(ctx, key)
}

// GetKeyBuilder returns the cache key builder
func (c *CacheService) GetKeyBuilder() *CacheKeyBuilder {
	return c.keyBuilder
}

// Close closes the cache connection
func (c *CacheService) Close() error {
	return c.cache.Close()
}

// TTL helper methods using configuration values
func (c *CacheService) GetAnimeDataTTL() time.Duration {
	return GetAnimeDataTTL(c.config)
}

func (c *CacheService) GetEpisodeTTL() time.Duration {
	return GetEpisodeTTL(c.config)
}

func (c *CacheService) GetSeasonTTL() time.Duration {
	return GetSeasonTTL(c.config)
}

func (c *CacheService) GetLockTTL() time.Duration {
	return GetLockTTL(c.config)
}