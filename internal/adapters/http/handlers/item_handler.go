package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/akonovalovdev/DDD_example/internal/ports/input"
)

// ItemHandler обрабатывает HTTP запросы для работы с предметами
type ItemHandler struct {
	service input.ItemService
	logger  *slog.Logger
}

// NewItemHandler создает новый ItemHandler
func NewItemHandler(service input.ItemService, logger *slog.Logger) *ItemHandler {
	return &ItemHandler{
		service: service,
		logger:  logger,
	}
}

// GetItems обрабатывает GET /items
func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	items, err := h.service.GetItems(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch items", h.logger)
		return
	}

	respondWithJSON(w, http.StatusOK, items, h.logger)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Warn("Failed to encode response", "error", err)
	}
}

func respondWithError(w http.ResponseWriter, status int, message string, logger *slog.Logger) {
	respondWithJSON(w, status, map[string]string{"error": message}, logger)
}
