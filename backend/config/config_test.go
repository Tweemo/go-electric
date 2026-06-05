package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("ENV", "production") // skip .env loading
	t.Setenv("HOST", "")
	t.Setenv("PORT", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	cfg := Load()

	if cfg.Host != "localhost" {
		t.Errorf("Host = %q, want localhost", cfg.Host)
	}
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.Addr() != "localhost:8080" {
		t.Errorf("Addr() = %q, want localhost:8080", cfg.Addr())
	}
	if len(cfg.CORSOrigins) != 2 {
		t.Errorf("CORSOrigins = %v, want 2 dev defaults", cfg.CORSOrigins)
	}
	if cfg.MaxUploadBytes != MaxUploadBytes {
		t.Errorf("MaxUploadBytes = %d, want %d", cfg.MaxUploadBytes, MaxUploadBytes)
	}
}

func TestLoadEnvOverride(t *testing.T) {
	t.Setenv("ENV", "production")
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("PORT", "9999")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://a.com, https://b.com ")

	cfg := Load()

	if cfg.Addr() != "0.0.0.0:9999" {
		t.Errorf("Addr() = %q, want 0.0.0.0:9999", cfg.Addr())
	}
	want := []string{"https://a.com", "https://b.com"}
	if len(cfg.CORSOrigins) != len(want) {
		t.Fatalf("CORSOrigins = %v, want %v", cfg.CORSOrigins, want)
	}
	for i, o := range want {
		if cfg.CORSOrigins[i] != o {
			t.Errorf("CORSOrigins[%d] = %q, want %q", i, cfg.CORSOrigins[i], o)
		}
	}
}
