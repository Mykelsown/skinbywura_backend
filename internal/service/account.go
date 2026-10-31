package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/Mykelsown/skinbywura_backend.git/internal/mailer"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/Mykelsown/skinbywura_backend.git/internal/types"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const resetTTL = time.Hour

// ErrInvalidContact is returned when a contact submission is missing required fields.
var ErrInvalidContact = errors.New("name, a valid email and a message are required")

// UpdateProfile changes the signed-in user's display name.
func (s *Service) UpdateProfile(ctx context.Context, user types.User, req types.UpdateProfileRequest) (types.User, error) {
	name := clean(req.Name, 80)
	if name == "" {
		return types.User{}, ErrInvalidInput
	}
	if err := s.store.UpdateUserName(ctx, user.ID, name); err != nil {
		return types.User{}, err
	}
	user.Name = name
	return user, nil
}

// ChangePassword verifies the current password, then sets the new one and signs out other devices.
func (s *Service) ChangePassword(ctx context.Context, user types.User, currentToken string, req types.ChangePasswordRequest) error {
	if !validPassword(req.NewPassword) {
		return ErrWeakPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return ErrWrongPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.store.UpdatePasswordKeepingSession(ctx, user.ID, string(hash), currentToken)
}

// RequestPasswordReset emails a single-use link if the address belongs to an account.
// It always returns nil for unknown addresses so the endpoint cannot be used to discover who has an account.
func (s *Service) RequestPasswordReset(ctx context.Context, req types.ForgotPasswordRequest) error {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if !validEmail(email) {
		return nil
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	token, err := randomToken()
	if err != nil {
		return err
	}
	if err := s.store.CreatePasswordReset(ctx, user.ID, hashToken(token), time.Now().Add(resetTTL)); err != nil {
		return err
	}

	link := s.opts.FrontendURL + "/reset-password?token=" + url.QueryEscape(token)
	body := fmt.Sprintf("Hi %s,\n\nUse this link to choose a new password. It works once and expires in one hour:\n\n%s\n\nIf you did not ask for this, you can ignore this email.\n\n%s",
		user.Name, link, s.opts.StoreName)

	s.sendAsync(mailer.Message{To: user.Email, Subject: "Reset your " + s.opts.StoreName + " password", Text: body})
	return nil
}

// ResetPassword completes a reset. Every existing session for the account is revoked.
func (s *Service) ResetPassword(ctx context.Context, req types.ResetPasswordRequest) error {
	if !validPassword(req.NewPassword) {
		return ErrWeakPassword
	}
	token := strings.TrimSpace(req.Token)
	if token == "" {
		return ErrResetInvalid
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.store.ConsumePasswordReset(ctx, hashToken(token), string(hash)); err != nil {
		if errors.Is(err, store.ErrResetInvalid) {
			return ErrResetInvalid
		}
		return err
	}
	return nil
}

// SubmitContact stores the message first, then alerts the store owner by email.
func (s *Service) SubmitContact(ctx context.Context, req types.ContactRequest) error {
	name := clean(req.Name, 80)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	topic := clean(req.Topic, 60)
	message := clean(req.Message, 4000)
	if name == "" || !validEmail(email) || message == "" {
		return ErrInvalidContact
	}
	if topic == "" {
		topic = "General"
	}

	if err := s.store.CreateContactMessage(ctx, name, email, topic, message); err != nil {
		return err
	}

	if s.opts.StoreEmail != "" {
		s.sendAsync(mailer.Message{
			To:      s.opts.StoreEmail,
			Subject: fmt.Sprintf("[%s] %s from %s", s.opts.StoreName, topic, name),
			Text:    fmt.Sprintf("From: %s <%s>\nTopic: %s\n\n%s", name, email, topic, message),
			ReplyTo: email,
		})
	}
	return nil
}

// CleanupLoop deletes expired sessions and reset tokens every hour until ctx ends.
func (s *Service) CleanupLoop(ctx context.Context) {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := s.store.DeleteExpired(ctx); err != nil {
				slog.Warn("cleanup failed", "err", err)
			}
		}
	}
}

// sendAsync sends mail without blocking the request. Failures are logged, never shown to shoppers,
// because the underlying action (order, message) has already been saved.
func (s *Service) sendAsync(m mailer.Message) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := s.opts.Mailer.Send(ctx, m); err != nil {
			slog.Error("email failed", "to", m.To, "subject", m.Subject, "err", err)
		}
	}()
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken stores only a SHA-256 of the reset token, so a database leak cannot be used to take over accounts.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
