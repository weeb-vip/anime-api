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

// OptimizedCacheService wraps CompressedCacheService with field optimization
type OptimizedCacheService struct {
	*CompressedCacheService
	excludeFields map[string][]string // Fields to exclude per type
}

// NewOptimizedCacheService creates a cache service with compression and field optimization
func NewOptimizedCacheService(cache Cache, cfg config.RedisConfig) *OptimizedCacheService {
	excludeFields := map[string][]string{
		"AnimeEpisode": {"Synopsis", "synopsis"},                    // Exclude episode synopses
		"Anime":       {"Synopsis", "synopsis", "TitleSynonyms", "title_synonyms", "Genres", "genres", "Licensors", "licensors", "Broadcast", "broadcast"},           // Exclude heavy text fields
		"Episode":     {"Synopsis", "synopsis"},                    // Exclude episode synopses
	}

	return &OptimizedCacheService{
		CompressedCacheService: NewCompressedCacheService(cache, cfg),
		excludeFields:         excludeFields,
	}
}

// SetJSONOptimized stores JSON data with field optimization for cache efficiency
func (o *OptimizedCacheService) SetJSONOptimized(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "OptimizedCache.SetJSONOptimized",
		trace.WithAttributes(
			attribute.String("cache.operation", "set_optimized"),
			attribute.String("cache.key", key),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	// Optimize the value before caching
	optimizedValue := o.optimizeForCache(value)

	// Count optimizations
	originalFields := o.countFields(value)
	optimizedFields := o.countFields(optimizedValue)
	fieldsRemoved := originalFields - optimizedFields

	span.SetAttributes(
		attribute.Int("cache.original_fields", originalFields),
		attribute.Int("cache.optimized_fields", optimizedFields),
		attribute.Int("cache.fields_removed", fieldsRemoved),
		attribute.Bool("cache.optimized", fieldsRemoved > 0),
	)

	return o.CompressedCacheService.SetJSON(ctx, key, optimizedValue, ttl)
}

// SetJSON automatically uses field optimization for all cache writes
func (o *OptimizedCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return o.SetJSONOptimized(ctx, key, value, ttl)
}

// optimizeForCache removes heavy fields from objects before caching
func (o *OptimizedCacheService) optimizeForCache(value interface{}) interface{} {
	return o.optimizeValue(reflect.ValueOf(value))
}

// optimizeValue recursively optimizes reflect.Value objects
func (o *OptimizedCacheService) optimizeValue(val reflect.Value) interface{} {
	// Handle pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		optimizedElem := o.optimizeValue(val.Elem())
		// Create new pointer to optimized value
		newPtr := reflect.New(reflect.TypeOf(optimizedElem))
		newPtr.Elem().Set(reflect.ValueOf(optimizedElem))
		return newPtr.Interface()
	}

	// Handle slices
	if val.Kind() == reflect.Slice {
		if val.IsNil() {
			return nil
		}

		optimizedSlice := reflect.MakeSlice(val.Type(), val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			optimizedItem := o.optimizeValue(val.Index(i))
			if optimizedItem != nil {
				optimizedSlice.Index(i).Set(reflect.ValueOf(optimizedItem))
			}
		}
		return optimizedSlice.Interface()
	}

	// Handle structs
	if val.Kind() == reflect.Struct {
		return o.optimizeStruct(val)
	}

	// For other types, return as-is
	return val.Interface()
}

// optimizeStruct removes heavy fields from struct values
func (o *OptimizedCacheService) optimizeStruct(val reflect.Value) interface{} {
	typ := val.Type()
	typeName := typ.Name()

	// Check if this type has fields to exclude
	excludeFields, hasExclusions := o.excludeFields[typeName]
	if !hasExclusions {
		// No optimizations for this type, but still process nested structs
		return o.optimizeNestedStructs(val)
	}

	// Create a new struct with optimized fields
	newVal := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip excluded fields
		if o.shouldExcludeField(fieldType.Name, excludeFields) {
			continue
		}

		// Recursively optimize the field value
		if newVal.Field(i).CanSet() {
			optimizedField := o.optimizeValue(field)
			if optimizedField != nil {
				newVal.Field(i).Set(reflect.ValueOf(optimizedField))
			}
		}
	}

	return newVal.Interface()
}

// optimizeNestedStructs processes nested structs even when parent doesn't need optimization
func (o *OptimizedCacheService) optimizeNestedStructs(val reflect.Value) interface{} {
	typ := val.Type()
	newVal := reflect.New(typ).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if newVal.Field(i).CanSet() {
			optimizedField := o.optimizeValue(field)
			if optimizedField != nil {
				newVal.Field(i).Set(reflect.ValueOf(optimizedField))
			}
		}
	}

	return newVal.Interface()
}

// shouldExcludeField checks if a field should be excluded from cache
func (o *OptimizedCacheService) shouldExcludeField(fieldName string, excludeFields []string) bool {
	for _, excludeField := range excludeFields {
		if fieldName == excludeField {
			return true
		}
	}
	return false
}

// countFields counts the total number of fields in a nested structure
func (o *OptimizedCacheService) countFields(value interface{}) int {
	return o.countFieldsInValue(reflect.ValueOf(value))
}

// countFieldsInValue recursively counts fields in reflect.Value
func (o *OptimizedCacheService) countFieldsInValue(val reflect.Value) int {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0
		}
		return o.countFieldsInValue(val.Elem())
	}

	if val.Kind() == reflect.Slice {
		count := 0
		for i := 0; i < val.Len(); i++ {
			count += o.countFieldsInValue(val.Index(i))
		}
		return count
	}

	if val.Kind() == reflect.Struct {
		count := val.NumField()
		for i := 0; i < val.NumField(); i++ {
			count += o.countFieldsInValue(val.Field(i))
		}
		return count
	}

	return 1 // Primitive field
}