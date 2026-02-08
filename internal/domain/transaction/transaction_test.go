package transaction

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestNewTransaction(t *testing.T) {
	userID := int64(1)
	amount := decimal.NewFromFloat(100.00)
	balanceBefore := decimal.NewFromFloat(1000.00)
	balanceAfter := decimal.NewFromFloat(900.00)
	description := "test transaction"

	tx := NewTransaction(userID, amount, balanceBefore, balanceAfter, description)

	if tx.ID == uuid.Nil {
		t.Error("expected non-nil UUID")
	}

	if tx.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, tx.UserID)
	}

	if !tx.Amount.Equal(amount) {
		t.Errorf("expected amount %s, got %s", amount.String(), tx.Amount.String())
	}

	if !tx.BalanceBefore.Equal(balanceBefore) {
		t.Errorf("expected balance before %s, got %s", balanceBefore.String(), tx.BalanceBefore.String())
	}

	if !tx.BalanceAfter.Equal(balanceAfter) {
		t.Errorf("expected balance after %s, got %s", balanceAfter.String(), tx.BalanceAfter.String())
	}

	if tx.Description != description {
		t.Errorf("expected description %s, got %s", description, tx.Description)
	}

	if tx.CreatedAt.IsZero() {
		t.Error("expected non-zero created at")
	}

	// Проверяем что время создания примерно сейчас (в пределах секунды)
	if time.Since(tx.CreatedAt) > time.Second {
		t.Error("created at should be approximately now")
	}
}

func TestNewWithdrawTransaction(t *testing.T) {
	userID := int64(1)
	amount := decimal.NewFromFloat(100.00)
	balanceBefore := decimal.NewFromFloat(1000.00)
	balanceAfter := decimal.NewFromFloat(900.00)

	tx := NewWithdrawTransaction(userID, amount, balanceBefore, balanceAfter)

	if tx.Description != "withdraw" {
		t.Errorf("expected description 'withdraw', got %s", tx.Description)
	}

	if tx.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, tx.UserID)
	}

	if !tx.Amount.Equal(amount) {
		t.Errorf("expected amount %s, got %s", amount.String(), tx.Amount.String())
	}
}

func TestTransaction_UniqueID(t *testing.T) {
	tx1 := NewTransaction(1, decimal.NewFromFloat(100), decimal.NewFromFloat(1000), decimal.NewFromFloat(900), "test1")
	tx2 := NewTransaction(1, decimal.NewFromFloat(100), decimal.NewFromFloat(1000), decimal.NewFromFloat(900), "test2")

	if tx1.ID == tx2.ID {
		t.Error("transactions should have unique IDs")
	}
}
