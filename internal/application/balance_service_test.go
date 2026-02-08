package application

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/transaction"
	"github.com/akonovalovdev/DDD_example/internal/domain/user"
)

type MockUserRepository struct {
	user       *user.User
	getUserErr error
	updateErr  error
	beginTxErr error
}

func (m *MockUserRepository) GetByID(_ context.Context, _ int64) (*user.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.user, nil
}

func (m *MockUserRepository) GetByIDForUpdate(_ context.Context, _ *sql.Tx, _ int64) (*user.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return user.NewUser(m.user.ID, m.user.Balance), nil
}

func (m *MockUserRepository) UpdateBalance(_ context.Context, _ *sql.Tx, _ int64, _ decimal.Decimal) error {
	return m.updateErr
}

func (m *MockUserRepository) BeginTx(_ context.Context) (*sql.Tx, error) {
	if m.beginTxErr != nil {
		return nil, m.beginTxErr
	}
	return nil, nil
}

type MockTransactionRepository struct {
	savedTransaction *transaction.Transaction
	saveErr          error
}

func (m *MockTransactionRepository) Save(_ context.Context, _ *sql.Tx, t *transaction.Transaction) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedTransaction = t
	return nil
}

func (m *MockTransactionRepository) GetByUserID(_ context.Context, _ int64, _, _ int) ([]*transaction.Transaction, error) {
	return nil, nil
}

func TestBalanceService_GetBalance_Success(t *testing.T) {
	expectedBalance := decimal.NewFromFloat(500.00)
	userRepo := &MockUserRepository{
		user: user.NewUser(1, expectedBalance),
	}
	txRepo := &MockTransactionRepository{}

	service := NewBalanceService(userRepo, txRepo)

	balance, err := service.GetBalance(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !balance.Equal(expectedBalance) {
		t.Errorf("expected balance %s, got %s", expectedBalance.String(), balance.String())
	}
}

func TestBalanceService_GetBalance_UserNotFound(t *testing.T) {
	userRepo := &MockUserRepository{
		getUserErr: user.ErrUserNotFound,
	}
	txRepo := &MockTransactionRepository{}

	service := NewBalanceService(userRepo, txRepo)

	_, err := service.GetBalance(context.Background(), 999)

	if !errors.Is(err, user.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestBalanceService_WithdrawBalance_BeginTxError(t *testing.T) {
	expectedErr := errors.New("connection failed")
	userRepo := &MockUserRepository{
		user:       user.NewUser(1, decimal.NewFromFloat(1000.00)),
		beginTxErr: expectedErr,
	}
	txRepo := &MockTransactionRepository{}

	service := NewBalanceService(userRepo, txRepo)

	_, err := service.WithdrawBalance(context.Background(), 1, decimal.NewFromFloat(100.00))

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
