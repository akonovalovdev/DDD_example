package input

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/transaction"
)

// WithdrawResult содержит результат операции списания
type WithdrawResult struct {
	Transaction   *transaction.Transaction
	BalanceBefore decimal.Decimal
	BalanceAfter  decimal.Decimal
}

// BalanceService определяет интерфейс сервиса для работы с балансом
type BalanceService interface {
	// WithdrawBalance списывает средства с баланса пользователя
	WithdrawBalance(ctx context.Context, userID int64, amount decimal.Decimal) (*WithdrawResult, error)

	// GetBalance возвращает текущий баланс пользователя
	GetBalance(ctx context.Context, userID int64) (decimal.Decimal, error)
}
