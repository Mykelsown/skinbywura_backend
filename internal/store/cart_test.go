package store

import (
	"testing"

	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

func TestMergeCartItems(t *testing.T) {
	in := []types.CartItem{
		{ProductID: 1, Qty: 2},
		{ID: 2, Qty: 1},        // falls back to ID when ProductID is zero
		{ProductID: 1, Qty: 3}, // duplicate merges
		{ProductID: 3, Qty: 0}, // dropped
		{ProductID: -4, Qty: 1},
		{ProductID: 5, Qty: 500}, // capped
	}
	got := mergeCartItems(in)
	want := []cartLine{{1, 5}, {2, 1}, {5, MaxLineQty}}
	if len(got) != len(want) {
		t.Fatalf("got %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d: got %#v want %#v", i, got[i], want[i])
		}
	}
}
