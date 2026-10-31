package service

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

const (
	minPasswordLen = 8
	maxPasswordLen = 72 // bcrypt silently ignores bytes beyond 72, so reject instead of truncating.
)

// validEmail accepts a plain address like a@b.com and rejects display-name forms.
func validEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email && strings.Contains(email, ".")
}

// validPassword enforces length bounds measured in bytes, matching bcrypt's limit.
func validPassword(p string) bool {
	return len(p) >= minPasswordLen && len(p) <= maxPasswordLen
}

// clean trims whitespace and caps the length in characters.
func clean(s string, max int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) > max {
		r := []rune(s)
		s = string(r[:max])
	}
	return s
}
