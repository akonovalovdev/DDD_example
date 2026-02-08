package item

import "errors"

var (
	// ErrItemNotFound возвращается когда предмет не найден
	ErrItemNotFound = errors.New("item not found")

	// ErrFetchFailed возвращается когда не удалось получить предметы
	ErrFetchFailed = errors.New("failed to fetch items")

	// ErrEmptyResponse возвращается когда API вернул пустой ответ
	ErrEmptyResponse = errors.New("empty response from API")
)
