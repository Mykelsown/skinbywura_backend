// Package mailer sends transactional email (password resets, order confirmations, contact alerts).
// With no API key configured it logs instead of sending, so local development never needs email setup.
package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// Message is one outgoing email.
type Message struct {
	To      string
	Subject string
	Text    string
	ReplyTo string
}

// Sender delivers a message. Implementations must be safe for concurrent use.
type Sender interface {
	Send(ctx context.Context, m Message) error
}

// LogSender prints messages to the log. Used when no email provider is configured.
type LogSender struct{}

// Send logs the email, including its body so reset links are usable in development.
func (LogSender) Send(_ context.Context, m Message) error {
	slog.Info("email (not sent, RESEND_API_KEY unset)", "to", m.To, "subject", m.Subject, "body", m.Text)
	return nil
}

// Resend sends mail through the Resend HTTP API using only the standard library.
type Resend struct {
	APIKey string
	From   string
	Client *http.Client
	URL    string
}

// NewResend builds a Resend sender with sane timeouts.
func NewResend(apiKey, from string) *Resend {
	return &Resend{
		APIKey: apiKey,
		From:   from,
		Client: &http.Client{Timeout: 10 * time.Second},
		URL:    "https://api.resend.com/emails",
	}
}

// Send posts the message to Resend and surfaces any non-2xx response as an error.
func (r *Resend) Send(ctx context.Context, m Message) error {
	payload := map[string]any{
		"from":    r.From,
		"to":      []string{m.To},
		"subject": m.Subject,
		"text":    m.Text,
	}
	if m.ReplyTo != "" {
		payload["reply_to"] = m.ReplyTo
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("resend returned %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// New picks the real sender when a key is present and the logger otherwise.
func New(apiKey, from string) Sender {
	if apiKey == "" {
		return LogSender{}
	}
	return NewResend(apiKey, from)
}
