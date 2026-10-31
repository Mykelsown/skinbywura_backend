package service

import (
	"strings"
	"testing"
)

func TestValidEmail(t *testing.T) {
	good := []string{"a@b.com", "first.last@shop.co.uk"}
	bad := []string{"", "nope", "a@b", "Name <a@b.com>", "a@b.com, c@d.com"}
	for _, e := range good {
		if !validEmail(e) {
			t.Errorf("%q should be valid", e)
		}
	}
	for _, e := range bad {
		if validEmail(e) {
			t.Errorf("%q should be invalid", e)
		}
	}
}

func TestValidPassword(t *testing.T) {
	if validPassword("short") {
		t.Error("short password accepted")
	}
	if !validPassword("longenough1") {
		t.Error("good password rejected")
	}
	if validPassword(strings.Repeat("a", 73)) {
		t.Error("over-72-byte password accepted")
	}
}

func TestCleanCapsLength(t *testing.T) {
	if got := clean("  hello  ", 10); got != "hello" {
		t.Errorf("got %q", got)
	}
	if got := clean("abcdef", 3); got != "abc" {
		t.Errorf("got %q", got)
	}
}
