package input

import (
	"context"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
)

// ItemService определяет интерфейс сервиса для работы с предметами
type ItemService interface {
	// GetItems возвращает список предметов с минимальными ценами
	GetItems(ctx context.Context) ([]*item.Item, error)
}
