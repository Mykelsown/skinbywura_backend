package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Mykelsown/skinbywura_backend.git/internal/mailer"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
)

var (
	ErrCartEmpty       = store.ErrCartEmpty
	ErrInvalidShipping = errors.New("name, phone, address, city and state are required for delivery")
)

// CreateOrder creates a user order from the current authenticated cart and notifies customer and store.
func (s *Service) CreateOrder(ctx context.Context, user types.User, req types.CheckoutRequest) (types.Order, error) {
	ship := types.CheckoutRequest{
		Name:    clean(req.Name, 80),
		Phone:   clean(req.Phone, 30),
		Address: clean(req.Address, 300),
		City:    clean(req.City, 80),
		State:   clean(req.State, 80),
		Notes:   clean(req.Notes, 500),
	}
	if ship.Name == "" || len(ship.Phone) < 7 || ship.Address == "" || ship.City == "" || ship.State == "" {
		return types.Order{}, ErrInvalidShipping
	}

	order, err := s.store.CreateOrder(ctx, user.ID, ship)
	if err != nil {
		if errors.Is(err, store.ErrCartEmpty) {
			return types.Order{}, ErrCartEmpty
		}
		return types.Order{}, err
	}

	s.notifyOrder(user, order)
	return order, nil
}

// ListOrdersByUser returns all orders for the given user.
func (s *Service) ListOrdersByUser(ctx context.Context, userID int64) ([]types.Order, error) {
	return s.store.ListOrdersByUser(ctx, userID)
}

func (s *Service) notifyOrder(user types.User, o types.Order) {
	summary := orderSummary(o)

	customer := fmt.Sprintf("Hi %s,\n\nThanks for your order with %s. We have received it and will confirm shortly.\n\n%s\n%s\n",
		user.Name, s.opts.StoreName, summary, deliveryBlock(o))
	if s.opts.PaymentNote != "" {
		customer += "\nPayment: " + s.opts.PaymentNote + "\n"
	}
	s.sendAsync(mailer.Message{To: user.Email, Subject: fmt.Sprintf("%s order #%d received", s.opts.StoreName, o.ID), Text: customer})

	if s.opts.StoreEmail != "" {
		alert := fmt.Sprintf("New order #%d from %s <%s>\n\n%s\n%s", o.ID, user.Name, user.Email, summary, deliveryBlock(o))
		s.sendAsync(mailer.Message{To: s.opts.StoreEmail, Subject: fmt.Sprintf("New order #%d (%s)", o.ID, formatNaira(o.Total)), Text: alert, ReplyTo: user.Email})
	}
}

func orderSummary(o types.Order) string {
	var b strings.Builder
	for _, it := range o.Items {
		fmt.Fprintf(&b, "- %d x %s  %s\n", it.Qty, it.Name, formatNaira(it.Price*it.Qty))
	}
	fmt.Fprintf(&b, "Total: %s\n", formatNaira(o.Total))
	return b.String()
}

func deliveryBlock(o types.Order) string {
	s := fmt.Sprintf("\nDeliver to:\n%s\n%s\n%s, %s\nPhone: %s\n", o.ShipName, o.ShipAddress, o.ShipCity, o.ShipState, o.ShipPhone)
	if o.Notes != "" {
		s += "Notes: " + o.Notes + "\n"
	}
	return s
}

// formatNaira renders whole naira with thousands separators, e.g. 12500 -> "NGN 12,500".
func formatNaira(n int) string {
	s := fmt.Sprintf("%d", n)
	var out []byte
	for i, c := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	return "NGN " + string(out)
}
