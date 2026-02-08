package output

import (
	"context"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
)

// ItemFetcher определяет интерфейс для получения предметов из внешнего источника
type ItemFetcher interface {
	// FetchItems получает список предметов из внешнего API
	FetchItems(ctx context.Context) ([]*item.Item, error)
}
