package cache

import (
	"context"
	"reflect"
	"time"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UltraOptimizedCacheService provides maximum cache optimization by separating large collections
type UltraOptimizedCacheService struct {
	*CompressedCacheService
	maxEpisodesInCache int    // Maximum episodes to include in anime cache
	excludeFields      map[string][]string
}

// NewUltraOptimizedCacheService creates a cache service with maximum optimization
func NewUltraOptimizedCacheService(cache Cache, cfg config.RedisConfig) *UltraOptimizedCacheService {
	excludeFields := map[string][]string{
		"AnimeEpisode": {"Synopsis", "synopsis", "TitleJp", "title_jp"},
		"Anime": {
			"Synopsis", "synopsis",
			"TitleSynonyms", "title_synonyms",
			"Genres", "genres",
			"Licensors", "licensors",
			"Broadcast", "broadcast",
			"Source", "source",
			"Studios", "studios",  // Keep only essential fields
		},
		"Episode": {"Synopsis", "synopsis", "TitleJp", "title_jp"},
	}

	return &UltraOptimizedCacheService{
		CompressedCacheService: NewCompressedCacheService(cache, cfg),
		maxEpisodesInCache:     0, // Don't cache episodes with anime at all
		excludeFields:         excludeFields,
	}
}

// SetJSON with ultra optimization - removes episodes entirely from anime cache
func (u *UltraOptimizedCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "UltraOptimizedCache.SetJSON",
		trace.WithAttributes(
			attribute.String("cache.operation", "set_ultra_optimized"),
			attribute.String("cache.key", key),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	// Ultra optimize the value
	optimizedValue, episodesSeparated := u.ultraOptimize(ctx, key, value, ttl)

	span.SetAttributes(
		attribute.Bool("cache.episodes_separated", episodesSeparated),
		attribute.Int("cache.max_episodes_in_cache", u.maxEpisodesInCache),
	)

	// Store the optimized value
	return u.CompressedCacheService.SetJSON(ctx, key, optimizedValue, ttl)
}

// ultraOptimize performs maximum optimization including episode separation
func (u *UltraOptimizedCacheService) ultraOptimize(ctx context.Context, key string, value interface{}, ttl time.Duration) (interface{}, bool) {
	// Handle nil interface values
	if value == nil {
		return value, false
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check for zero/invalid value after pointer dereferencing
	if !val.IsValid() {
		return value, false
	}

	// Special handling for time types - don't optimize them
	if val.Type().String() == "time.Time" {
		return value, false
	}

	if val.Kind() == reflect.Slice {
		// For slices (like anime lists), optimize each item
		return u.optimizeSlice(ctx, key, val, ttl), false
	}

	if val.Kind() == reflect.Struct {
		return u.optimizeStruct(ctx, key, val, ttl)
	}

	return value, false
}

// optimizeSlice optimizes each item in a slice
func (u *UltraOptimizedCacheService) optimizeSlice(ctx context.Context, key string, val reflect.Value, ttl time.Duration) interface{} {
	// Check if slice value is valid
	if !val.IsValid() || val.Kind() != reflect.Slice {
		// Return original interface if val is not valid to avoid nil returns
		if val.IsValid() {
			return val.Interface()
		}
		return nil
	}

	optimizedSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Cap())

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if !item.IsValid() {
			// Copy original item if it's invalid
			optimizedSlice.Index(i).Set(item)
			continue
		}
		optimizedItem, _ := u.ultraOptimize(ctx, key, item.Interface(), ttl)
		if optimizedItem != nil {
			optimizedValue := reflect.ValueOf(optimizedItem)
			if optimizedValue.IsValid() && optimizedValue.Type().AssignableTo(optimizedSlice.Index(i).Type()) {
				optimizedSlice.Index(i).Set(optimizedValue)
			} else {
				// Fallback to original item if optimization failed
				optimizedSlice.Index(i).Set(item)
			}
		} else {
			// Fallback to original item if optimization returned nil
			optimizedSlice.Index(i).Set(item)
		}
	}

	return optimizedSlice.Interface()
}

// optimizeStruct optimizes struct fields and handles episode separation
func (u *UltraOptimizedCacheService) optimizeStruct(ctx context.Context, key string, val reflect.Value, ttl time.Duration) (interface{}, bool) {
	// Check if struct value is valid
	if !val.IsValid() || val.Kind() != reflect.Struct {
		return val.Interface(), false
	}

	typ := val.Type()
	typeName := typ.Name()
	episodesSeparated := false

	// Create optimized struct
	newVal := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Name

		// Skip invalid fields
		if !field.IsValid() {
			continue
		}

		// Handle episodes field specially
		if fieldName == "AnimeEpisodes" || fieldName == "Episodes" {
			if field.Kind() == reflect.Slice && field.Len() > u.maxEpisodesInCache {
				// Separate episodes into their own cache entries
				u.cacheEpisodesSeparately(ctx, key, field, ttl)
				episodesSeparated = true
				// Don't include episodes in main cache entry
				continue
			}
		}

		// Skip excluded fields
		if u.shouldExcludeField(fieldName, typeName) {
			continue
		}

		// Recursively optimize nested fields
		if newVal.Field(i).CanSet() {
			optimizedField, _ := u.ultraOptimize(ctx, key, field.Interface(), ttl)
			if optimizedField != nil {
				optimizedVal := reflect.ValueOf(optimizedField)
				// Validate the optimized value
				if !optimizedVal.IsValid() {
					// Fallback to original field value
					newVal.Field(i).Set(field)
					continue
				}
				fieldType := newVal.Field(i).Type()

				// Handle type mismatches, especially for time types
				if optimizedVal.Type() != fieldType {
					// Skip optimization if types don't match, keep original value
					newVal.Field(i).Set(field)
					continue
				}

				newVal.Field(i).Set(optimizedVal)
			} else {
				// Fallback to original field value if optimization returned nil
				newVal.Field(i).Set(field)
			}
		}
	}

	return newVal.Interface(), episodesSeparated
}

// cacheEpisodesSeparately stores episodes in separate cache entries
func (u *UltraOptimizedCacheService) cacheEpisodesSeparately(ctx context.Context, baseKey string, episodesField reflect.Value, ttl time.Duration) {
	if episodesField.Kind() != reflect.Slice {
		return
	}

	// Store episodes in separate cache key
	episodesCacheKey := baseKey + ":episodes"

	// Optimize episodes before storing
	optimizedEpisodes := u.optimizeSlice(ctx, episodesCacheKey, episodesField, ttl)

	// Store episodes separately (async to not slow down main cache operation)
	go func() {
		_ = u.CompressedCacheService.SetJSON(context.Background(), episodesCacheKey, optimizedEpisodes, ttl)
	}()
}

// shouldExcludeField checks if a field should be excluded from cache
func (u *UltraOptimizedCacheService) shouldExcludeField(fieldName, typeName string) bool {
	excludeFields, hasExclusions := u.excludeFields[typeName]
	if !hasExclusions {
		return false
	}

	for _, excludeField := range excludeFields {
		if fieldName == excludeField {
			return true
		}
	}
	return false
}