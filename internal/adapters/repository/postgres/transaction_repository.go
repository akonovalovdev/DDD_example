package postgres

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/transaction"
)

// TransactionRepository реализует репозиторий транзакций для PostgreSQL
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository создает новый экземпляр TransactionRepository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Save сохраняет транзакцию
func (r *TransactionRepository) Save(ctx context.Context, tx *sql.Tx, t *transaction.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, amount, balance_before, balance_after, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		t.ID,
		t.UserID,
		t.Amount.String(),
		t.BalanceBefore.String(),
		t.BalanceAfter.String(),
		t.Description,
		t.CreatedAt,
	)

	return err
}

// GetByUserID возвращает список транзакций пользователя
func (r *TransactionRepository) GetByUserID(ctx context.Context, userID int64, limit, offset int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, user_id, amount, balance_before, balance_after, description, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	for rows.Next() {
		var t transaction.Transaction
		var amount, balanceBefore, balanceAfter string

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&amount,
			&balanceBefore,
			&balanceAfter,
			&t.Description,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		t.Amount, _ = decimal.NewFromString(amount)
		t.BalanceBefore, _ = decimal.NewFromString(balanceBefore)
		t.BalanceAfter, _ = decimal.NewFromString(balanceAfter)

		transactions = append(transactions, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
