package output

import (
	"context"
	"time"
)

// Cache определяет интерфейс для кэширования данных
type Cache interface {
	// Get получает значение из кэша
	Get(ctx context.Context, key string) (interface{}, bool)

	// Set устанавливает значение в кэш с указанным TTL
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration)

	// Delete удаляет значение из кэша
	Delete(ctx context.Context, key string)

	// Clear очищает весь кэш
	Clear(ctx context.Context)
}
