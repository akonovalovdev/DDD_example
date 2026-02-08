package user

import "github.com/shopspring/decimal"

// CanWithdraw проверяет, может ли пользователь снять указанную сумму
func (u *User) CanWithdraw(amount decimal.Decimal) bool {
	if amount.LessThanOrEqual(decimal.Zero) {
		return false
	}
	return u.Balance.GreaterThanOrEqual(amount)
}

// Withdraw списывает сумму с баланса пользователя
// Возвращает баланс до операции и ошибку если операция невозможна
func (u *User) Withdraw(amount decimal.Decimal) (balanceBefore decimal.Decimal, err error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, ErrInvalidAmount
	}

	if u.Balance.LessThan(amount) {
		return decimal.Zero, ErrInsufficientBalance
	}

	balanceBefore = u.Balance
	u.Balance = u.Balance.Sub(amount)

	return balanceBefore, nil
}

// GetBalance возвращает текущий баланс
func (u *User) GetBalance() decimal.Decimal {
	return u.Balance
}
