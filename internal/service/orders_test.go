package service

import (
	"strings"
	"testing"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

func TestFormatNaira(t *testing.T) {
	cases := map[int]string{0: "NGN 0", 999: "NGN 999", 1000: "NGN 1,000", 12500: "NGN 12,500", 1234567: "NGN 1,234,567"}
	for in, want := range cases {
		if got := formatNaira(in); got != want {
			t.Errorf("formatNaira(%d)=%q want %q", in, got, want)
		}
	}
}

func TestOrderSummaryListsItemsAndTotal(t *testing.T) {
	o := types.Order{Total: 20000, Items: []types.OrderItem{{Name: "Serum", Price: 10000, Qty: 2}}}
	got := orderSummary(o)
	if !strings.Contains(got, "2 x Serum") || !strings.Contains(got, "NGN 20,000") {
		t.Fatalf("unexpected summary: %s", got)
	}
}

func TestHashTokenIsStableAndNotPlain(t *testing.T) {
	if hashToken("abc") != hashToken("abc") || hashToken("abc") == "abc" {
		t.Fatal("hashToken must be deterministic and not the identity")
	}
}
