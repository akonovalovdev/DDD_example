package application

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/transaction"
	"github.com/akonovalovdev/DDD_example/internal/ports/input"
	"github.com/akonovalovdev/DDD_example/internal/ports/output"
)

// BalanceServiceImpl реализует сервис для работы с балансом
type BalanceServiceImpl struct {
	userRepo        output.UserRepository
	transactionRepo output.TransactionRepository
}

// NewBalanceService создает новый экземпляр BalanceService
func NewBalanceService(
	userRepo output.UserRepository,
	transactionRepo output.TransactionRepository,
) *BalanceServiceImpl {
	return &BalanceServiceImpl{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// WithdrawBalance списывает средства с баланса пользователя
func (s *BalanceServiceImpl) WithdrawBalance(
	ctx context.Context,
	userID int64,
	amount decimal.Decimal,
) (*input.WithdrawResult, error) {
	// 1. Начинаем транзакцию БД
	tx, err := s.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 2. Получаем пользователя с блокировкой (SELECT ... FOR UPDATE)
	user, err := s.userRepo.GetByIDForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// 3. Выполняем domain логику — списание
	balanceBefore, err := user.Withdraw(amount)
	if err != nil {
		return nil, err
	}

	// 4. Создаем запись истории транзакции
	txRecord := transaction.NewWithdrawTransaction(
		userID,
		amount,
		balanceBefore,
		user.Balance,
	)

	// 5. Сохраняем транзакцию в историю
	if err = s.transactionRepo.Save(ctx, tx, txRecord); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %w", err)
	}

	// 6. Обновляем баланс пользователя
	if err = s.userRepo.UpdateBalance(ctx, tx, userID, user.Balance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// 7. Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &input.WithdrawResult{
		Transaction:   txRecord,
		BalanceBefore: balanceBefore,
		BalanceAfter:  user.Balance,
	}, nil
}

// GetBalance возвращает текущий баланс пользователя
func (s *BalanceServiceImpl) GetBalance(ctx context.Context, userID int64) (decimal.Decimal, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return decimal.Zero, err
	}
	return user.Balance, nil
}
