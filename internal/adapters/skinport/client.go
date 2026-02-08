package skinport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/shopspring/decimal"

	"github.com/akonovalovdev/DDD_example/internal/domain/item"
)

// SkinportItem представляет предмет из API Skinport
type SkinportItem struct {
	MarketHashName string   `json:"market_hash_name"`
	MinPrice       *float64 `json:"min_price"`
	SuggestedPrice *float64 `json:"suggested_price"`
}

// Client реализует клиент для Skinport API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый клиент Skinport API
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchItems получает список предметов из Skinport API
func (c *Client) FetchItems(ctx context.Context) ([]*item.Item, error) {
	// Делаем два запроса параллельно: tradable и non-tradable
	tradableCh := make(chan fetchResult)
	nonTradableCh := make(chan fetchResult)

	go func() {
		items, err := c.fetchItems(ctx, true)
		tradableCh <- fetchResult{items: items, err: err}
	}()

	go func() {
		items, err := c.fetchItems(ctx, false)
		nonTradableCh <- fetchResult{items: items, err: err}
	}()

	tradableResult := <-tradableCh
	nonTradableResult := <-nonTradableCh

	if tradableResult.err != nil {
		return nil, fmt.Errorf("failed to fetch tradable items: %w", tradableResult.err)
	}
	if nonTradableResult.err != nil {
		return nil, fmt.Errorf("failed to fetch non-tradable items: %w", nonTradableResult.err)
	}

	// Объединяем результаты
	return mergeItems(tradableResult.items, nonTradableResult.items), nil
}

type fetchResult struct {
	items map[string]*SkinportItem
	err   error
}

func (c *Client) fetchItems(ctx context.Context, tradable bool) (map[string]*SkinportItem, error) {
	url := fmt.Sprintf("%s/items?app_id=730&currency=USD&tradable=%t", c.baseURL, tradable)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Skinport API требует поддержку Brotli компрессии
	req.Header.Set("Accept-Encoding", "br")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декомпрессия Brotli если сервер вернул сжатый ответ
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "br" {
		reader = brotli.NewReader(resp.Body)
	}

	var items []SkinportItem
	if err := json.NewDecoder(reader).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Преобразуем в map для быстрого поиска
	result := make(map[string]*SkinportItem, len(items))
	for i := range items {
		result[items[i].MarketHashName] = &items[i]
	}

	return result, nil
}

func mergeItems(tradable, nonTradable map[string]*SkinportItem) []*item.Item {
	// Собираем все уникальные имена предметов
	allNames := make(map[string]struct{})
	for name := range tradable {
		allNames[name] = struct{}{}
	}
	for name := range nonTradable {
		allNames[name] = struct{}{}
	}

	result := make([]*item.Item, 0, len(allNames))

	for name := range allNames {
		tradablePrice := decimal.Zero
		nonTradablePrice := decimal.Zero

		if t, ok := tradable[name]; ok && t.MinPrice != nil {
			tradablePrice = decimal.NewFromFloat(*t.MinPrice)
		}

		if nt, ok := nonTradable[name]; ok && nt.MinPrice != nil {
			nonTradablePrice = decimal.NewFromFloat(*nt.MinPrice)
		}

		result = append(result, item.NewItem(name, tradablePrice, nonTradablePrice))
	}

	return result
}
