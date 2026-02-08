package cache

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryCache_SetAndGet(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Minute)

	val, ok := cache.Get(ctx, "key1")
	if !ok {
		t.Error("expected key to exist")
	}

	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}
}

func TestInMemoryCache_GetNonExistent(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	val, ok := cache.Get(ctx, "nonexistent")
	if ok {
		t.Error("expected key to not exist")
	}

	if val != nil {
		t.Errorf("expected nil, got %v", val)
	}
}

func TestInMemoryCache_Expiration(t *testing.T) {
	cache := NewInMemoryCache(time.Millisecond * 100)
	defer cache.Close()

	ctx := context.Background()

	// Устанавливаем с очень коротким TTL
	cache.Set(ctx, "expiring", "value", time.Millisecond*50)

	// Сразу должно быть доступно
	val, ok := cache.Get(ctx, "expiring")
	if !ok {
		t.Error("expected key to exist immediately after set")
	}
	if val != "value" {
		t.Errorf("expected value, got %v", val)
	}

	// Ждём истечения TTL
	time.Sleep(time.Millisecond * 60)

	// Теперь должно быть недоступно
	_, ok = cache.Get(ctx, "expiring")
	if ok {
		t.Error("expected key to be expired")
	}
}

func TestInMemoryCache_Delete(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	cache.Set(ctx, "to_delete", "value", time.Minute)

	// Проверяем что существует
	_, ok := cache.Get(ctx, "to_delete")
	if !ok {
		t.Error("expected key to exist before delete")
	}

	// Удаляем
	cache.Delete(ctx, "to_delete")

	// Проверяем что удалилось
	_, ok = cache.Get(ctx, "to_delete")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestInMemoryCache_Clear(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	cache.Set(ctx, "key1", "value1", time.Minute)
	cache.Set(ctx, "key2", "value2", time.Minute)
	cache.Set(ctx, "key3", "value3", time.Minute)

	if cache.Len() != 3 {
		t.Errorf("expected 3 items, got %d", cache.Len())
	}

	cache.Clear(ctx)

	if cache.Len() != 0 {
		t.Errorf("expected 0 items after clear, got %d", cache.Len())
	}
}

func TestInMemoryCache_Len(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	if cache.Len() != 0 {
		t.Errorf("expected 0 items initially, got %d", cache.Len())
	}

	cache.Set(ctx, "key1", "value1", time.Minute)
	if cache.Len() != 1 {
		t.Errorf("expected 1 item, got %d", cache.Len())
	}

	cache.Set(ctx, "key2", "value2", time.Minute)
	if cache.Len() != 2 {
		t.Errorf("expected 2 items, got %d", cache.Len())
	}
}

func TestInMemoryCache_OverwriteKey(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	cache.Set(ctx, "key", "value1", time.Minute)
	cache.Set(ctx, "key", "value2", time.Minute)

	val, ok := cache.Get(ctx, "key")
	if !ok {
		t.Error("expected key to exist")
	}

	if val != "value2" {
		t.Errorf("expected value2, got %v", val)
	}

	if cache.Len() != 1 {
		t.Errorf("expected 1 item (overwritten), got %d", cache.Len())
	}
}

func TestInMemoryCache_DifferentTypes(t *testing.T) {
	cache := NewInMemoryCache(time.Minute)
	defer cache.Close()

	ctx := context.Background()

	// Строка
	cache.Set(ctx, "string", "hello", time.Minute)
	if val, _ := cache.Get(ctx, "string"); val != "hello" {
		t.Errorf("expected hello, got %v", val)
	}

	// Число
	cache.Set(ctx, "int", 42, time.Minute)
	if val, _ := cache.Get(ctx, "int"); val != 42 {
		t.Errorf("expected 42, got %v", val)
	}

	// Срез
	slice := []string{"a", "b", "c"}
	cache.Set(ctx, "slice", slice, time.Minute)
	if val, _ := cache.Get(ctx, "slice"); len(val.([]string)) != 3 {
		t.Errorf("expected slice of 3, got %v", val)
	}

	// Структура
	type TestStruct struct {
		Name string
		Age  int
	}
	cache.Set(ctx, "struct", TestStruct{Name: "test", Age: 25}, time.Minute)
	if val, _ := cache.Get(ctx, "struct"); val.(TestStruct).Name != "test" {
		t.Errorf("expected test, got %v", val)
	}
}
