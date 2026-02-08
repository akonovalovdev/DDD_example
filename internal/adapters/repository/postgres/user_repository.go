package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/user"
)

// UserRepository реализует репозиторий пользователей для PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	query := `SELECT id, balance FROM users WHERE id = $1`

	var u user.User
	var balance string

	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	u.Balance, _ = decimal.NewFromString(balance)
	return &u, nil
}

// GetByIDForUpdate возвращает пользователя по ID с блокировкой для обновления
func (r *UserRepository) GetByIDForUpdate(ctx context.Context, tx *sql.Tx, id int64) (*user.User, error) {
	query := `SELECT id, balance FROM users WHERE id = $1 FOR UPDATE`

	var u user.User
	var balance string

	err := tx.QueryRowContext(ctx, query, id).Scan(&u.ID, &balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	u.Balance, _ = decimal.NewFromString(balance)
	return &u, nil
}

// UpdateBalance обновляет баланс пользователя
func (r *UserRepository) UpdateBalance(ctx context.Context, tx *sql.Tx, id int64, balance decimal.Decimal) error {
	query := `UPDATE users SET balance = $1 WHERE id = $2`

	result, err := tx.ExecContext(ctx, query, balance.String(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

// BeginTx начинает транзакцию
func (r *UserRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
}
