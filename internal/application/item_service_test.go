package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
)

type MockItemFetcher struct {
	items []*item.Item
	err   error
}

func (m *MockItemFetcher) FetchItems(ctx context.Context) ([]*item.Item, error) {
	return m.items, m.err
}

type MockCache struct {
	data map[string]interface{}
}

func NewMockCache() *MockCache {
	return &MockCache{data: make(map[string]interface{})}
}

func (m *MockCache) Get(ctx context.Context, key string) (interface{}, bool) {
	val, ok := m.data[key]
	return val, ok
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	m.data[key] = value
}

func (m *MockCache) Delete(ctx context.Context, key string) {
	delete(m.data, key)
}

func (m *MockCache) Clear(ctx context.Context) {
	m.data = make(map[string]interface{})
}

func TestItemService_GetItems_FromAPI(t *testing.T) {
	price1 := decimal.NewFromFloat(100)
	price2 := decimal.NewFromFloat(90)
	price3 := decimal.NewFromFloat(200)
	price4 := decimal.NewFromFloat(180)

	expectedItems := []*item.Item{
		{MarketHashName: "AK-47", TradableMinPrice: &price1, NonTradableMinPrice: &price2},
		{MarketHashName: "AWP", TradableMinPrice: &price3, NonTradableMinPrice: &price4},
	}

	fetcher := &MockItemFetcher{items: expectedItems}
	cache := NewMockCache()

	service := NewItemService(fetcher, cache, 5*time.Minute)

	items, err := service.GetItems(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != len(expectedItems) {
		t.Errorf("expected %d items, got %d", len(expectedItems), len(items))
	}

	cached, ok := cache.Get(context.Background(), "skinport:items")
	if !ok {
		t.Error("expected items to be cached")
	}

	cachedItems := cached.([]*item.Item)
	if len(cachedItems) != len(expectedItems) {
		t.Errorf("expected %d cached items, got %d", len(expectedItems), len(cachedItems))
	}
}

func TestItemService_GetItems_FromCache(t *testing.T) {
	price1 := decimal.NewFromFloat(100)
	price2 := decimal.NewFromFloat(90)

	cachedItems := []*item.Item{
		{MarketHashName: "Cached AK-47", TradableMinPrice: &price1, NonTradableMinPrice: &price2},
	}

	fetcher := &MockItemFetcher{
		err: errors.New("should not be called"),
	}

	cache := NewMockCache()
	cache.Set(context.Background(), "skinport:items", cachedItems, 5*time.Minute)

	service := NewItemService(fetcher, cache, 5*time.Minute)

	items, err := service.GetItems(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Errorf("expected 1 item from cache, got %d", len(items))
	}

	if items[0].MarketHashName != "Cached AK-47" {
		t.Errorf("expected cached item name, got %s", items[0].MarketHashName)
	}
}

func TestItemService_GetItems_FetchError(t *testing.T) {
	expectedError := errors.New("fetch failed")

	fetcher := &MockItemFetcher{err: expectedError}
	cache := NewMockCache()

	service := NewItemService(fetcher, cache, 5*time.Minute)

	_, err := service.GetItems(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
