package config

import "testing"

func TestValidateProductionRules(t *testing.T) {
	base := EnvData{Env: "production", DB: "postgres://x", CookieSecure: true, CookieSameSite: "lax", AllowedOrigins: []string{"https://shop.example.com"}}
	if err := base.Validate(); err != nil {
		t.Fatalf("valid production config rejected: %v", err)
	}

	bad := base
	bad.CookieSecure = false
	if bad.Validate() == nil {
		t.Fatal("insecure cookie must be rejected in production")
	}

	bad = base
	bad.AllowedOrigins = []string{"http://localhost:5173"}
	if bad.Validate() == nil {
		t.Fatal("localhost origin must be rejected in production")
	}

	bad = base
	bad.DB = ""
	if bad.Validate() == nil {
		t.Fatal("missing DATABASE_URL must be rejected")
	}
}

func TestSplitListTrimsSlashesAndBlanks(t *testing.T) {
	got := splitList(" https://a.com/ , ,https://b.com ")
	if len(got) != 2 || got[0] != "https://a.com" || got[1] != "https://b.com" {
		t.Fatalf("unexpected result: %#v", got)
	}
}
