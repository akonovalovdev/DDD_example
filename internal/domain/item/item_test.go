package item

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestItem_Creation(t *testing.T) {
	name := "AK-47 | Redline"
	tradablePrice := decimal.NewFromFloat(12.50)
	nonTradablePrice := decimal.NewFromFloat(10.20)

	item := &Item{
		MarketHashName:      name,
		TradableMinPrice:    &tradablePrice,
		NonTradableMinPrice: &nonTradablePrice,
	}

	if item.MarketHashName != name {
		t.Errorf("expected market hash name %s, got %s", name, item.MarketHashName)
	}

	if !item.TradableMinPrice.Equal(tradablePrice) {
		t.Errorf("expected tradable price %s, got %s", tradablePrice.String(), item.TradableMinPrice.String())
	}

	if !item.NonTradableMinPrice.Equal(nonTradablePrice) {
		t.Errorf("expected non-tradable price %s, got %s", nonTradablePrice.String(), item.NonTradableMinPrice.String())
	}
}

func TestItem_NilPrices(t *testing.T) {
	item := &Item{
		MarketHashName:      "Test Item",
		TradableMinPrice:    nil,
		NonTradableMinPrice: nil,
	}

	if item.TradableMinPrice != nil {
		t.Error("expected nil tradable price")
	}

	if item.NonTradableMinPrice != nil {
		t.Error("expected nil non-tradable price")
	}
}
