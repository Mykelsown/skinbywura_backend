// Package config loads the application's configuration values from environment variables.
// Every deploy-specific value lives here so the same binary runs locally and in production.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// EnvData contains the runtime values that are read from the environment.
type EnvData struct {
	Env            string        // Env is "development" or "production".
	Port           string        // Port is the TCP port used by the HTTP server.
	DB             string        // DB is the PostgreSQL connection string.
	SessionSecret  string        // SessionSecret is reserved for signing; sessions are random DB tokens today.
	AllowedOrigins []string      // AllowedOrigins are the browser origins allowed to call the API with cookies.
	FrontendURL    string        // FrontendURL is the public storefront URL, used in email links.
	CookieSecure   bool          // CookieSecure marks the session cookie HTTPS-only.
	CookieSameSite string        // CookieSameSite is "lax", "strict" or "none".
	AutoMigrate    bool          // AutoMigrate applies pending SQL migrations at startup.
	ResendAPIKey   string        // ResendAPIKey enables real email delivery when set.
	MailFrom       string        // MailFrom is the From header, e.g. "Skin By Wura <orders@yourdomain.com>".
	StoreEmail     string        // StoreEmail receives contact messages and new-order alerts.
	StoreName      string        // StoreName is used in email copy.
	PaymentNote    string        // PaymentNote is the payment instructions text added to order emails.
	RequestTimeout time.Duration // RequestTimeout caps how long a single request may run.
}

// IsProduction reports whether the app is running in production mode.
func (c EnvData) IsProduction() bool { return c.Env == "production" }

// Load reads the configured environment variables and returns them as a structured config object.
func Load() EnvData {
	env := strings.ToLower(getenv("APP_ENV", "development"))
	return EnvData{
		Env:            env,
		Port:           getenv("PORT", "8080"),
		DB:             os.Getenv("DATABASE_URL"),
		SessionSecret:  os.Getenv("SESSION_SECRET"),
		AllowedOrigins: splitList(getenv("ALLOWED_ORIGINS", "http://localhost:5173")),
		FrontendURL:    strings.TrimRight(getenv("FRONTEND_URL", "http://localhost:5173"), "/"),
		CookieSecure:   getbool("COOKIE_SECURE", env == "production"),
		CookieSameSite: strings.ToLower(getenv("COOKIE_SAMESITE", "lax")),
		AutoMigrate:    getbool("AUTO_MIGRATE", true),
		ResendAPIKey:   os.Getenv("RESEND_API_KEY"),
		MailFrom:       getenv("MAIL_FROM", "Skin By Wura <onboarding@resend.dev>"),
		StoreEmail:     os.Getenv("STORE_EMAIL"),
		StoreName:      getenv("STORE_NAME", "Skin By Wura"),
		PaymentNote:    os.Getenv("PAYMENT_NOTE"),
		RequestTimeout: 30 * time.Second,
	}
}

// Validate fails fast on configuration that would make the deploy unsafe or broken.
func (c EnvData) Validate() error {
	if c.DB == "" {
		return errors.New("DATABASE_URL is required")
	}
	switch c.CookieSameSite {
	case "lax", "strict", "none":
	default:
		return errors.New("COOKIE_SAMESITE must be lax, strict or none")
	}
	if c.CookieSameSite == "none" && !c.CookieSecure {
		return errors.New("COOKIE_SAMESITE=none requires COOKIE_SECURE=true")
	}
	if c.IsProduction() {
		if !c.CookieSecure {
			return errors.New("COOKIE_SECURE must be true in production")
		}
		for _, o := range c.AllowedOrigins {
			if o == "*" || strings.HasPrefix(o, "http://localhost") {
				return errors.New("ALLOWED_ORIGINS must list the real storefront origin in production")
			}
		}
	}
	return nil
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getbool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func splitList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimRight(strings.TrimSpace(p), "/"); p != "" {
			out = append(out, p)
		}
	}
	return out
}
