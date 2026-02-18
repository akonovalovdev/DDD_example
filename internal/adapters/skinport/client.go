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
Currency       string   `json:"currency"`
SuggestedPrice *float64 `json:"suggested_price"`
ItemPage       string   `json:"item_page"`
MarketPage     string   `json:"market_page"`
MinPrice       *float64 `json:"min_price"`
MaxPrice       *float64 `json:"max_price"`
MeanPrice      *float64 `json:"mean_price"`
Quantity       int      `json:"quantity"`
CreatedAt      int64    `json:"created_at"`
UpdatedAt      int64    `json:"updated_at"`
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
		var tradablePrice, nonTradablePrice *decimal.Decimal
		var suggestedPrice, maxPrice, meanPrice *decimal.Decimal
		var currency, itemPage, marketPage string
		var quantity int
		var createdAt, updatedAt int64

		// Берём данные из tradable
		if t, ok := tradable[name]; ok {
			currency = t.Currency
			itemPage = t.ItemPage
			marketPage = t.MarketPage
			quantity = t.Quantity
			createdAt = t.CreatedAt
			updatedAt = t.UpdatedAt

			if t.MinPrice != nil {
				p := decimal.NewFromFloat(*t.MinPrice)
				tradablePrice = &p
			}
			if t.SuggestedPrice != nil {
				p := decimal.NewFromFloat(*t.SuggestedPrice)
				suggestedPrice = &p
			}
			if t.MaxPrice != nil {
				p := decimal.NewFromFloat(*t.MaxPrice)
				maxPrice = &p
			}
			if t.MeanPrice != nil {
				p := decimal.NewFromFloat(*t.MeanPrice)
				meanPrice = &p
			}
		}

		// Берём non-tradable цену
		if nt, ok := nonTradable[name]; ok {
			if nt.MinPrice != nil {
				p := decimal.NewFromFloat(*nt.MinPrice)
				nonTradablePrice = &p
			}
			// Заполняем остальные поля если не были заполнены из tradable
			if currency == "" {
				currency = nt.Currency
				itemPage = nt.ItemPage
				marketPage = nt.MarketPage
				quantity = nt.Quantity
				createdAt = nt.CreatedAt
				updatedAt = nt.UpdatedAt
			}
		}

		result = append(result, &item.Item{
			MarketHashName:      name,
			Currency:            currency,
			SuggestedPrice:      suggestedPrice,
			ItemPage:            itemPage,
			MarketPage:          marketPage,
			TradableMinPrice:    tradablePrice,
			NonTradableMinPrice: nonTradablePrice,
			MaxPrice:            maxPrice,
			MeanPrice:           meanPrice,
			Quantity:            quantity,
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
		})
	}

	return result
}
