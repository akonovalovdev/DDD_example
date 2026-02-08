package user

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestUser_Withdraw_Success(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(1000.00))

	balanceBefore, err := user.Withdraw(decimal.NewFromFloat(100.00))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !balanceBefore.Equal(decimal.NewFromFloat(1000.00)) {
		t.Errorf("expected balance before 1000.00, got %s", balanceBefore.String())
	}

	if !user.Balance.Equal(decimal.NewFromFloat(900.00)) {
		t.Errorf("expected balance after 900.00, got %s", user.Balance.String())
	}
}

func TestUser_Withdraw_InsufficientBalance(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(50.00))

	_, err := user.Withdraw(decimal.NewFromFloat(100.00))

	if err != ErrInsufficientBalance {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}

	if !user.Balance.Equal(decimal.NewFromFloat(50.00)) {
		t.Errorf("balance should not change on failed withdraw, got %s", user.Balance.String())
	}
}

func TestUser_Withdraw_ZeroAmount(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(100.00))

	_, err := user.Withdraw(decimal.Zero)

	if err != ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestUser_Withdraw_NegativeAmount(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(100.00))

	_, err := user.Withdraw(decimal.NewFromFloat(-50.00))

	if err != ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestUser_Withdraw_ExactBalance(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(100.00))

	balanceBefore, err := user.Withdraw(decimal.NewFromFloat(100.00))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !balanceBefore.Equal(decimal.NewFromFloat(100.00)) {
		t.Errorf("expected balance before 100.00, got %s", balanceBefore.String())
	}

	if !user.Balance.Equal(decimal.Zero) {
		t.Errorf("expected balance after 0, got %s", user.Balance.String())
	}
}

func TestUser_CanWithdraw(t *testing.T) {
	tests := []struct {
		name     string
		balance  float64
		amount   float64
		expected bool
	}{
		{"sufficient balance", 100.00, 50.00, true},
		{"exact balance", 100.00, 100.00, true},
		{"insufficient balance", 50.00, 100.00, false},
		{"zero amount", 100.00, 0, false},
		{"negative amount", 100.00, -50.00, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(1, decimal.NewFromFloat(tt.balance))
			result := user.CanWithdraw(decimal.NewFromFloat(tt.amount))

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUser_GetBalance(t *testing.T) {
	user := NewUser(1, decimal.NewFromFloat(500.50))

	balance := user.GetBalance()

	if !balance.Equal(decimal.NewFromFloat(500.50)) {
		t.Errorf("expected 500.50, got %s", balance.String())
	}
}

func TestNewUser(t *testing.T) {
	user := NewUser(42, decimal.NewFromFloat(1234.56))

	if user.ID != 42 {
		t.Errorf("expected ID 42, got %d", user.ID)
	}

	if !user.Balance.Equal(decimal.NewFromFloat(1234.56)) {
		t.Errorf("expected balance 1234.56, got %s", user.Balance.String())
	}
}
