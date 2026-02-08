package output

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/user"
)

// UserRepository определяет интерфейс репозитория для работы с пользователями
type UserRepository interface {
	// GetByID возвращает пользователя по ID
	GetByID(ctx context.Context, id int64) (*user.User, error)

	// GetByIDForUpdate возвращает пользователя по ID с блокировкой для обновления
	GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*user.User, error)

	// UpdateBalance обновляет баланс пользователя
	UpdateBalance(ctx context.Context, tx *sql.Tx, id int64, balance decimal.Decimal) error

	// BeginTx начинает транзакцию
	BeginTx(ctx context.Context) (*sql.Tx, error)
}
