package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/user"
	"github.com/akonovalovdev/DDD_example/internal/ports/input"
)

// BalanceHandler обрабатывает HTTP запросы для работы с балансом
type BalanceHandler struct {
	service input.BalanceService
	logger  *slog.Logger
}

// NewBalanceHandler создает новый BalanceHandler
func NewBalanceHandler(service input.BalanceService, logger *slog.Logger) *BalanceHandler {
	return &BalanceHandler{
		service: service,
		logger:  logger,
	}
}

// WithdrawRequest представляет запрос на списание
type WithdrawRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

// WithdrawResponse представляет ответ на списание
type WithdrawResponse struct {
	Success       bool            `json:"success"`
	TransactionID string          `json:"transaction_id"`
	BalanceBefore decimal.Decimal `json:"balance_before"`
	BalanceAfter  decimal.Decimal `json:"balance_after"`
}

// Withdraw обрабатывает POST /users/{id}/withdraw
func (h *BalanceHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Извлекаем userID из URL
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "user id is required", h.logger)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id", h.logger)
		return
	}

	// Декодируем тело запроса
	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body", h.logger)
		return
	}

	// Валидация суммы
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		respondWithError(w, http.StatusBadRequest, "amount must be positive", h.logger)
		return
	}

	// Выполняем списание
	result, err := h.service.WithdrawBalance(ctx, userID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			respondWithError(w, http.StatusNotFound, "user not found", h.logger)
		case errors.Is(err, user.ErrInsufficientBalance):
			respondWithError(w, http.StatusBadRequest, "insufficient balance", h.logger)
		case errors.Is(err, user.ErrInvalidAmount):
			respondWithError(w, http.StatusBadRequest, "invalid amount", h.logger)
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error", h.logger)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, WithdrawResponse{
		Success:       true,
		TransactionID: result.Transaction.ID.String(),
		BalanceBefore: result.BalanceBefore,
		BalanceAfter:  result.BalanceAfter,
	}, h.logger)
}

// GetBalance обрабатывает GET /users/{id}/balance
func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "user id is required", h.logger)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user id", h.logger)
		return
	}

	balance, err := h.service.GetBalance(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, "user not found", h.logger)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "internal server error", h.logger)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"balance": balance,
	}, h.logger)
}
