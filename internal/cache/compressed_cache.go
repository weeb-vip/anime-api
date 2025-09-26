package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"github.com/goccy/go-json"
)

// CompressedCacheService wraps CacheService with compression for large values
type CompressedCacheService struct {
	*CacheService
	compressionThreshold int // Compress values larger than this (bytes)
}

// NewCompressedCacheService creates a cache service with compression and high-performance JSON
func NewCompressedCacheService(cache Cache, cfg config.RedisConfig) *CompressedCacheService {
	return &CompressedCacheService{
		CacheService:         NewCacheService(cache, cfg),
		compressionThreshold: 1024, // Compress anything > 1KB
	}
}

// GetJSON retrieves and decompresses JSON data from cache
func (c *CompressedCacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "CompressedCache.GetJSON",
		trace.WithAttributes(
			attribute.String("cache.operation", "get_compressed"),
			attribute.String("cache.key", key),
			attribute.String("cache.layer", "compressed"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	// Get raw data from cache (Redis operation)
	redisStartTime := time.Now()
	data, err := c.cache.Get(ctx, key)
	redisEndTime := time.Now()
	redisDuration := redisEndTime.Sub(redisStartTime)

	span.SetAttributes(
		attribute.Int64("cache.redis_duration_us", redisDuration.Microseconds()),
		attribute.Int64("cache.redis_duration_ms", redisDuration.Milliseconds()),
	)
	if err != nil {
		return err
	}

	span.SetAttributes(
		attribute.Int("cache.compressed_size_bytes", len(data)),
	)

	// Check if data is compressed (starts with gzip magic number)
	var jsonData []byte
	if len(data) >= 2 && data[0] == 0x1f && data[1] == 0x8b {
		// Data is compressed, decompress it
		decompressStartTime := time.Now()
		_, decompressSpan := tracer.Start(ctx, "CompressedCache.Decompress",
			trace.WithSpanKind(trace.SpanKindInternal),
		)

		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			decompressSpan.RecordError(err)
			decompressSpan.End()
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.Close()

		jsonData, err = io.ReadAll(reader)
		decompressEndTime := time.Now()
		decompressDuration := decompressEndTime.Sub(decompressStartTime)

		if err != nil {
			decompressSpan.RecordError(err)
			decompressSpan.SetAttributes(
				attribute.Int64("cache.decompress_duration_us", decompressDuration.Microseconds()),
				attribute.Int64("cache.decompress_duration_ms", decompressDuration.Milliseconds()),
			)
			decompressSpan.End()
			return fmt.Errorf("failed to decompress data: %w", err)
		}

		decompressSpan.SetAttributes(
			attribute.Int("cache.decompressed_size_bytes", len(jsonData)),
			attribute.Float64("cache.compression_ratio", float64(len(data))/float64(len(jsonData))),
			attribute.Int64("cache.decompress_duration_us", decompressDuration.Microseconds()),
			attribute.Int64("cache.decompress_duration_ms", decompressDuration.Milliseconds()),
		)
		decompressSpan.SetStatus(codes.Ok, "decompression successful")
		decompressSpan.End()

		span.SetAttributes(
			attribute.Bool("cache.was_compressed", true),
			attribute.Int("cache.decompressed_size_bytes", len(jsonData)),
		)
	} else {
		// Data is not compressed
		jsonData = data
		span.SetAttributes(attribute.Bool("cache.was_compressed", false))
	}

	// Unmarshal JSON (using goccy/go-json for maximum performance)
	unmarshalStartTime := time.Now()
	err = json.Unmarshal(jsonData, dest)
	unmarshalEndTime := time.Now()
	unmarshalDuration := unmarshalEndTime.Sub(unmarshalStartTime)

	span.SetAttributes(
		attribute.Int64("cache.unmarshal_duration_us", unmarshalDuration.Microseconds()),
		attribute.Int64("cache.unmarshal_duration_ms", unmarshalDuration.Milliseconds()),
		attribute.Int("cache.json_size_bytes", len(jsonData)),
	)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.String("cache.result", "unmarshal_error"))
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	totalDuration := time.Since(startTime)
	span.SetAttributes(
		attribute.String("cache.result", "hit"),
		attribute.Int64("cache.total_duration_us", totalDuration.Microseconds()),
		attribute.Int64("cache.total_duration_ms", totalDuration.Milliseconds()),
	)
	span.SetStatus(codes.Ok, "cache hit with decompression")

	return nil
}

// SetJSON compresses and stores JSON data in cache
func (c *CompressedCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Marshal to JSON (using goccy/go-json for maximum performance)
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Run compression and cache set asynchronously - fire and forget
	go func() {
		tracer := tracing.GetTracer(context.Background())
		asyncCtx, span := tracer.Start(context.Background(), "CompressedCache.SetJSON",
			trace.WithAttributes(
				attribute.String("cache.operation", "set_compressed"),
				attribute.String("cache.key", key),
			),
			trace.WithSpanKind(trace.SpanKindInternal),
			tracing.GetEnvironmentAttribute(),
		)
		defer span.End()

		startTime := time.Now()
		originalSize := len(jsonData)
		span.SetAttributes(attribute.Int("cache.original_size_bytes", originalSize))

		var dataToStore []byte
		shouldCompress := originalSize > c.compressionThreshold

		if shouldCompress {
			// Compress the data
			_, compressSpan := tracer.Start(asyncCtx, "CompressedCache.Compress",
				trace.WithSpanKind(trace.SpanKindInternal),
			)

			var buf bytes.Buffer
			writer := gzip.NewWriter(&buf)
			_, err = writer.Write(jsonData)
			if err != nil {
				compressSpan.RecordError(err)
				compressSpan.End()
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return
			}
			writer.Close()

			dataToStore = buf.Bytes()
			compressedSize := len(dataToStore)

			compressSpan.SetAttributes(
				attribute.Int("cache.compressed_size_bytes", compressedSize),
				attribute.Float64("cache.compression_ratio", float64(compressedSize)/float64(originalSize)),
				attribute.Float64("cache.space_saved_percent", (1.0-float64(compressedSize)/float64(originalSize))*100),
			)
			compressSpan.SetStatus(codes.Ok, "compression successful")
			compressSpan.End()

			span.SetAttributes(
				attribute.Bool("cache.was_compressed", true),
				attribute.Int("cache.compressed_size_bytes", compressedSize),
				attribute.Float64("cache.compression_ratio", float64(compressedSize)/float64(originalSize)),
			)
		} else {
			// Don't compress small values
			dataToStore = jsonData
			span.SetAttributes(attribute.Bool("cache.was_compressed", false))
		}

		// Store in cache
		err = c.cache.Set(asyncCtx, key, dataToStore, ttl)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return
		}

		span.SetAttributes(
			attribute.String("cache.result", "success"),
			attribute.Int64("cache.duration_ms", time.Since(startTime).Milliseconds()),
		)
		span.SetStatus(codes.Ok, "cache set with compression")
	}()

	return nil
}