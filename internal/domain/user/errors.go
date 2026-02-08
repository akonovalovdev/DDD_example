package user

import "errors"

var (
	// ErrUserNotFound возвращается когда пользователь не найден
	ErrUserNotFound = errors.New("user not found")

	// ErrInsufficientBalance возвращается когда недостаточно средств на балансе
	ErrInsufficientBalance = errors.New("insufficient balance")

	// ErrInvalidAmount возвращается когда сумма некорректна (отрицательная или ноль)
	ErrInvalidAmount = errors.New("invalid amount: must be positive")

	// ErrUserAlreadyExists возвращается когда пользователь уже существует
	ErrUserAlreadyExists = errors.New("user already exists")
)
