package cache

import (
	"context"
	"sync"
	"time"
)

// item представляет элемент кэша
type item struct {
	value      interface{}
	expiration time.Time
}

// isExpired проверяет, истек ли срок действия элемента
func (i *item) isExpired() bool {
	return time.Now().After(i.expiration)
}

// InMemoryCache реализует in-memory кэш с TTL
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]*item

	// Для автоматической очистки
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

// NewInMemoryCache создает новый in-memory кэш
func NewInMemoryCache(cleanupInterval time.Duration) *InMemoryCache {
	c := &InMemoryCache{
		items:           make(map[string]*item),
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}

	// Запускаем горутину для очистки устаревших элементов
	go c.cleanup()

	return c
}

// Get получает значение из кэша
func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if item.isExpired() {
		return nil, false
	}

	return item.value, true
}

// Set устанавливает значение в кэш с указанным TTL
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &item{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete удаляет значение из кэша
func (c *InMemoryCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear очищает весь кэш
func (c *InMemoryCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*item)
}

// Close останавливает горутину очистки
func (c *InMemoryCache) Close() {
	close(c.stopCleanup)
}

// cleanup периодически очищает устаревшие элементы
func (c *InMemoryCache) cleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCleanup:
			return
		}
	}
}

// deleteExpired удаляет все устаревшие элементы
func (c *InMemoryCache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.isExpired() {
			delete(c.items, key)
		}
	}
}

// Len возвращает количество элементов в кэше
func (c *InMemoryCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}
