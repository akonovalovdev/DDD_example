package item

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewItem(t *testing.T) {
	name := "AK-47 | Redline"
	tradablePrice := decimal.NewFromFloat(12.50)
	nonTradablePrice := decimal.NewFromFloat(10.20)

	item := NewItem(name, tradablePrice, nonTradablePrice)

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

func TestNewItem_ZeroPrices(t *testing.T) {
	item := NewItem("Test Item", decimal.Zero, decimal.Zero)

	if !item.TradableMinPrice.Equal(decimal.Zero) {
		t.Error("expected zero tradable price")
	}

	if !item.NonTradableMinPrice.Equal(decimal.Zero) {
		t.Error("expected zero non-tradable price")
	}
}
