package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction представляет историю операции со счетом пользователя
type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	UserID        int64           `json:"user_id"`
	Amount        decimal.Decimal `json:"amount"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
	Description   string          `json:"description,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
}

// NewTransaction создает новую транзакцию
func NewTransaction(
	userID int64,
	amount decimal.Decimal,
	balanceBefore decimal.Decimal,
	balanceAfter decimal.Decimal,
	description string,
) *Transaction {
	return &Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		CreatedAt:     time.Now().UTC(),
	}
}

// NewWithdrawTransaction создает транзакцию списания
func NewWithdrawTransaction(
	userID int64,
	amount decimal.Decimal,
	balanceBefore decimal.Decimal,
	balanceAfter decimal.Decimal,
) *Transaction {
	return NewTransaction(
		userID,
		amount,
		balanceBefore,
		balanceAfter,
		"withdraw",
	)
}
