package item

import "github.com/shopspring/decimal"

// Item представляет предмет из Skinport с минимальными ценами
type Item struct {
	MarketHashName      string           `json:"market_hash_name"`
	Currency            string           `json:"currency"`
	SuggestedPrice      *decimal.Decimal `json:"suggested_price,omitempty"`
	ItemPage            string           `json:"item_page"`
	MarketPage          string           `json:"market_page"`
	TradableMinPrice    *decimal.Decimal `json:"tradable_min_price,omitempty"`
	NonTradableMinPrice *decimal.Decimal `json:"non_tradable_min_price,omitempty"`
	MaxPrice            *decimal.Decimal `json:"max_price,omitempty"`
	MeanPrice           *decimal.Decimal `json:"mean_price,omitempty"`
	Quantity            int              `json:"quantity"`
	CreatedAt           int64            `json:"created_at"`
	UpdatedAt           int64            `json:"updated_at"`
}
