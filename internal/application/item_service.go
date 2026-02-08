package application

import (
	"context"
	"time"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
	"github.com/akonovalovdev/DDD_example/internal/ports/output"
)

const itemsCacheKey = "skinport:items"

// ItemServiceImpl реализует сервис для работы с предметами
type ItemServiceImpl struct {
	fetcher  output.ItemFetcher
	cache    output.Cache
	cacheTTL time.Duration
}

// NewItemService создает новый экземпляр ItemService
func NewItemService(
	fetcher output.ItemFetcher,
	cache output.Cache,
	cacheTTL time.Duration,
) *ItemServiceImpl {
	return &ItemServiceImpl{
		fetcher:  fetcher,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

// GetItems возвращает список предметов с минимальными ценами
func (s *ItemServiceImpl) GetItems(ctx context.Context) ([]*item.Item, error) {
	// 1. Проверяем кэш
	if cached, ok := s.cache.Get(ctx, itemsCacheKey); ok {
		if items, ok := cached.([]*item.Item); ok {
			return items, nil
		}
	}

	// 2. Если в кэше нет — запрашиваем из API
	items, err := s.fetcher.FetchItems(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Сохраняем в кэш
	s.cache.Set(ctx, itemsCacheKey, items, s.cacheTTL)

	return items, nil
}
