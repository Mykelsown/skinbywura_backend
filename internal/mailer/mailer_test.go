package mailer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResendSendsExpectedPayload(t *testing.T) {
	var got map[string]any
	var auth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		_ = json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := NewResend("key123", "Shop <a@b.com>")
	s.URL = srv.URL
	err := s.Send(context.Background(), Message{To: "x@y.com", Subject: "Hi", Text: "Body", ReplyTo: "r@y.com"})
	if err != nil {
		t.Fatal(err)
	}
	if auth != "Bearer key123" {
		t.Fatalf("bad auth header %q", auth)
	}
	if got["subject"] != "Hi" || got["reply_to"] != "r@y.com" {
		t.Fatalf("bad payload %#v", got)
	}
}

func TestResendReportsProviderErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}))
	defer srv.Close()
	s := NewResend("k", "f")
	s.URL = srv.URL
	if err := s.Send(context.Background(), Message{To: "x@y.com"}); err == nil {
		t.Fatal("expected error on 422")
	}
}

func TestNewFallsBackToLogSender(t *testing.T) {
	if _, ok := New("", "f").(LogSender); !ok {
		t.Fatal("empty key must give LogSender")
	}
}
