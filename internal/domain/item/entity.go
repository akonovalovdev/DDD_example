package item

import "github.com/shopspring/decimal"

// Item представляет предмет из Skinport с минимальными ценами
type Item struct {
	MarketHashName      string          `json:"market_hash_name"`
	TradableMinPrice    decimal.Decimal `json:"tradable_min_price"`
	NonTradableMinPrice decimal.Decimal `json:"non_tradable_min_price"`
}

// NewItem создает новый Item
func NewItem(marketHashName string, tradableMinPrice, nonTradableMinPrice decimal.Decimal) *Item {
	return &Item{
		MarketHashName:      marketHashName,
		TradableMinPrice:    tradableMinPrice,
		NonTradableMinPrice: nonTradableMinPrice,
	}
}
