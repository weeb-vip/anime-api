package resolvers

import (
	"context"
	"testing"
	"time"

	"github.com/weeb-vip/anime-api/internal/cache"
)

// MockCacheService implements CacheServiceInterface for testing
type MockCacheService struct {
	storage map[string]interface{}
	getCalls int
	setCalls int
}

func NewMockCacheService() *MockCacheService {
	return &MockCacheService{
		storage: make(map[string]interface{}),
	}
}

func (m *MockCacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	m.getCalls++
	if value, exists := m.storage[key]; exists {
		// Simple copy for testing - in real implementation this would be JSON unmarshaling
		switch value.(type) {
		case []interface{}:
			// Copy slice to dest (simplified for testing)
			return nil
		}
		return nil
	}
	return cache.ErrCacheMiss
}

func (m *MockCacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.setCalls++
	m.storage[key] = value
	return nil
}

func (m *MockCacheService) GetKeyBuilder() *cache.CacheKeyBuilder {
	return cache.NewCacheKeyBuilder("test")
}

func (m *MockCacheService) GetCurrentlyAiringTTL() time.Duration {
	return 5 * time.Minute
}

func TestCacheKeyGeneration(t *testing.T) {
	mockCache := NewMockCacheService()
	keyBuilder := mockCache.GetKeyBuilder()

	// Test cache key generation with different parameters
	key1 := keyBuilder.CurrentlyAiring(10, "", "", 0)
	key2 := keyBuilder.CurrentlyAiring(5, "", "", 0)
	key3 := keyBuilder.CurrentlyAiring(10, "2023-12-15", "", 0)

	if key1 == key2 {
		t.Error("Expected different cache keys for different limits")
	}

	if key1 == key3 {
		t.Error("Expected different cache keys for different start dates")
	}

	t.Logf("Generated cache keys:")
	t.Logf("  Default (limit=10): %s", key1)
	t.Logf("  Limit=5: %s", key2)
	t.Logf("  With start date: %s", key3)
}

func TestCacheServiceInterface(t *testing.T) {
	mockCache := NewMockCacheService()

	// Test that the mock implements the interface
	var _ CacheServiceInterface = mockCache

	// Test basic caching operations
	key := "test-key"
	value := []interface{}{"test", "data"}

	// Set value
	err := mockCache.SetJSON(context.Background(), key, value, 5*time.Minute)
	if err != nil {
		t.Errorf("SetJSON failed: %v", err)
	}

	if mockCache.setCalls != 1 {
		t.Errorf("Expected 1 set call, got %d", mockCache.setCalls)
	}

	// Get value
	var result []interface{}
	err = mockCache.GetJSON(context.Background(), key, &result)
	if err != nil {
		t.Errorf("GetJSON failed: %v", err)
	}

	if mockCache.getCalls != 1 {
		t.Errorf("Expected 1 get call, got %d", mockCache.getCalls)
	}

	// Test cache miss
	err = mockCache.GetJSON(context.Background(), "non-existent", &result)
	if err != cache.ErrCacheMiss {
		t.Errorf("Expected cache miss error, got: %v", err)
	}

	t.Log("Cache service interface test passed")
}