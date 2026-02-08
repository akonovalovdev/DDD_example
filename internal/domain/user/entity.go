package user

import "github.com/shopspring/decimal"

// User представляет пользователя системы
type User struct {
	ID      int64           `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

// NewUser создает нового пользователя
func NewUser(id int64, balance decimal.Decimal) *User {
	return &User{
		ID:      id,
		Balance: balance,
	}
}
