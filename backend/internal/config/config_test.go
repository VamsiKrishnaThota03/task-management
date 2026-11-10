package config

import (
	"strings"
	"testing"
)

// LoadConfig reads the process environment, so each case clears the variables
// it cares about. t.Setenv restores them automatically when the test ends.
func setEnv(t *testing.T, vars map[string]string) {
	t.Helper()
	for _, key := range []string{"DB_URI", "JWT_SECRET", "GEMINI_KEY", "GEMINI_MODEL", "HOST", "PORT", "CORS_ORIGINS"} {
		t.Setenv(key, "")
	}
	for key, value := range vars {
		t.Setenv(key, value)
	}
}

const validSecret = "0123456789abcdef0123456789abcdef" // exactly 32 chars

func TestLoadConfigRequiresDBURIAndJWTSecret(t *testing.T) {
	setEnv(t, map[string]string{})

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected an error when DB_URI and JWT_SECRET are both missing")
	}
	// The message must name both, so a misconfigured deploy is self-explaining.
	for _, want := range []string{"DB_URI", "JWT_SECRET"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should name %s, got: %v", want, err)
		}
	}
}

func TestLoadConfigRejectsAShortJWTSecret(t *testing.T) {
	// A short signing secret is brute-forceable, so this is a hard failure
	// rather than a warning.
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": "too-short",
	})

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected an error for a JWT_SECRET under 32 characters")
	}
	if !strings.Contains(err.Error(), "32") {
		t.Errorf("error should mention the 32-character minimum, got: %v", err)
	}
}

func TestLoadConfigAcceptsASecretAtExactlyTheMinimum(t *testing.T) {
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
	})

	if len(validSecret) != minJWTSecretLength {
		t.Fatalf("test fixture is wrong: want %d chars, got %d", minJWTSecretLength, len(validSecret))
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("a secret of exactly the minimum length should be accepted: %v", err)
	}
	if cfg.JWTSecret != validSecret {
		t.Errorf("JWTSecret = %q, want %q", cfg.JWTSecret, validSecret)
	}
}

func TestLoadConfigDefaultsToLoopback(t *testing.T) {
	// This is the security-relevant default: an unset HOST must not expose the
	// dev server to the local network.
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Host != "127.0.0.1" {
		t.Errorf("Host = %q, want 127.0.0.1 — an unset HOST must not bind all interfaces", cfg.Host)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
}

func TestLoadConfigHonoursAnExplicitHost(t *testing.T) {
	// Container hosting needs an all-interfaces bind; it just has to be asked for.
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
		"HOST":       "0.0.0.0",
		"PORT":       "9999",
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Host != "0.0.0.0" {
		t.Errorf("Host = %q, want 0.0.0.0", cfg.Host)
	}
	if cfg.Port != "9999" {
		t.Errorf("Port = %q, want 9999", cfg.Port)
	}
}

func TestLoadConfigParsesCORSOrigins(t *testing.T) {
	setEnv(t, map[string]string{
		"DB_URI":       "file::memory:",
		"JWT_SECRET":   validSecret,
		"CORS_ORIGINS": " https://app.example.com/ ,http://localhost:3000, ,https://other.example/ ",
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	want := []string{"https://app.example.com", "http://localhost:3000", "https://other.example"}
	if len(cfg.AllowedOrigins) != len(want) {
		t.Fatalf("AllowedOrigins = %v, want %v", cfg.AllowedOrigins, want)
	}
	for i, origin := range want {
		if cfg.AllowedOrigins[i] != origin {
			t.Errorf("AllowedOrigins[%d] = %q, want %q (trailing slashes and blanks must be stripped)",
				i, cfg.AllowedOrigins[i], origin)
		}
	}
}

func TestLoadConfigDefaultsCORSToLocalhost(t *testing.T) {
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "http://localhost:3000" {
		t.Errorf("AllowedOrigins = %v, want [http://localhost:3000]", cfg.AllowedOrigins)
	}
}

func TestLoadConfigStripsQuotesFromGeminiKey(t *testing.T) {
	// A key pasted into .env with surrounding quotes is a common mistake; the
	// loader strips them rather than sending a quoted key upstream.
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
		"GEMINI_KEY": `"a-quoted-key"`,
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.GeminiKey != "a-quoted-key" {
		t.Errorf("GeminiKey = %q, want a-quoted-key", cfg.GeminiKey)
	}
}

func TestLoadConfigSucceedsWithoutAGeminiKey(t *testing.T) {
	// AI suggestions are optional: a missing key logs a notice, it does not
	// stop the service from starting.
	setEnv(t, map[string]string{
		"DB_URI":     "file::memory:",
		"JWT_SECRET": validSecret,
	})

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("a missing GEMINI_KEY must not be fatal: %v", err)
	}
	if cfg.GeminiKey != "" {
		t.Errorf("GeminiKey = %q, want empty", cfg.GeminiKey)
	}
	if cfg.GeminiModel != "gemini-2.5-flash" {
		t.Errorf("GeminiModel = %q, want the default gemini-2.5-flash", cfg.GeminiModel)
	}
}

func TestLoadConfigTrimsWhitespaceAroundRequiredValues(t *testing.T) {
	// A trailing newline from a secret manager must not count as a value, nor
	// inflate the length check past the minimum.
	setEnv(t, map[string]string{
		"DB_URI":     "   ",
		"JWT_SECRET": "  " + validSecret + "  ",
	})

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("whitespace-only DB_URI should count as missing")
	}
	if !strings.Contains(err.Error(), "DB_URI") {
		t.Errorf("error should name DB_URI, got: %v", err)
	}
}
