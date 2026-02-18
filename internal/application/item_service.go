package application

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
	"github.com/akonovalovdev/DDD_example/internal/ports/output"
)

const itemsCacheKey = "skinport:items"

// ItemServiceImpl реализует сервис для работы с предметами
type ItemServiceImpl struct {
	fetcher  output.ItemFetcher
	cache    output.Cache
	cacheTTL time.Duration
	sfGroup  singleflight.Group // защита от thundering herd
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

// WarmUp прогревает кеш при запуске приложения
func (s *ItemServiceImpl) WarmUp(ctx context.Context) error {
	_, err := s.GetItems(ctx)
	return err
}

// GetItems возвращает список предметов с минимальными ценами
func (s *ItemServiceImpl) GetItems(ctx context.Context) ([]*item.Item, error) {
	// 1. Проверяем кэш
	if cached, ok := s.cache.Get(ctx, itemsCacheKey); ok {
		if items, ok := cached.([]*item.Item); ok {
			return items, nil
		}
	}

	// 2. Singleflight — дедупликация параллельных запросов
	result, err, _ := s.sfGroup.Do(itemsCacheKey, func() (interface{}, error) {
		// Повторная проверка кеша (мог заполниться пока ждали)
		if cached, ok := s.cache.Get(ctx, itemsCacheKey); ok {
			if items, ok := cached.([]*item.Item); ok {
				return items, nil
			}
		}

		// Запрашиваем из API
		items, err := s.fetcher.FetchItems(ctx)
		if err != nil {
			return nil, err
		}

		// Сохраняем в кэш
		s.cache.Set(ctx, itemsCacheKey, items, s.cacheTTL)

		return items, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*item.Item), nil
}
