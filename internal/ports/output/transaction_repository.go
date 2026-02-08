package output

import (
	"context"
	"database/sql"

	"github.com/akonovalovdev/DDD_example/internal/domain/transaction"
)

// TransactionRepository определяет интерфейс репозитория для работы с транзакциями
type TransactionRepository interface {
	// Save сохраняет транзакцию
	Save(ctx context.Context, tx *sql.Tx, t *transaction.Transaction) error

	// GetByUserID возвращает список транзакций пользователя
	GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*transaction.Transaction, error)
}
