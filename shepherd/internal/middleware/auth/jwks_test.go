package auth

import (
	"testing"
	"time"
)

// TestMaxAge проверяет разбор Cache-Control, по которому кешируется JWKS.
// Ошибка здесь либо заставляет ходить в Keycloak на каждый промах локального
// кеша, либо, наоборот, держит устаревшие ключи после ротации.
func TestMaxAge(t *testing.T) {
	tests := []struct {
		name         string
		cacheControl string
		want         time.Duration
	}{
		{name: "заголовок отсутствует", cacheControl: "", want: defaultTTL},
		{name: "только max-age", cacheControl: "max-age=300", want: 300 * time.Second},
		{name: "среди других директив", cacheControl: "public, max-age=600, must-revalidate", want: 600 * time.Second},
		{name: "с пробелами вокруг", cacheControl: "  max-age=42  ", want: 42 * time.Second},
		// max-age=0 значит «не кешировать», но fromKeycloak всё равно кладёт
		// значение в Redis, поэтому нулевой TTL заменяется на defaultTTL:
		// нулевой TTL в Redis означает «хранить вечно», а это хуже.
		{name: "нулевой max-age даёт значение по умолчанию", cacheControl: "max-age=0", want: defaultTTL},
		{name: "отрицательный max-age", cacheControl: "max-age=-5", want: defaultTTL},
		{name: "нечисловой max-age", cacheControl: "max-age=abc", want: defaultTTL},
		{name: "без max-age", cacheControl: "no-cache, no-store", want: defaultTTL},
		// s-maxage — отдельная директива для прокси, путать её с max-age нельзя.
		{name: "s-maxage не считается за max-age", cacheControl: "s-maxage=100", want: defaultTTL},
		{name: "первый корректный max-age", cacheControl: "max-age=30, max-age=90", want: 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxAge(tt.cacheControl); got != tt.want {
				t.Errorf("maxAge(%q) = %v, ожидалось %v", tt.cacheControl, got, tt.want)
			}
		})
	}
}

// TestNewJWKSFetcher фиксирует вывод адреса JWKS и ключа кеша из issuer.
// Realm берётся последним сегментом — при опечатке fetcher молча пойдёт не в
// тот realm и все токены станут невалидными.
func TestNewJWKSFetcher(t *testing.T) {
	tests := []struct {
		name         string
		issuerURL    string
		wantJWKSURL  string
		wantCacheKey string
	}{
		{
			name:         "обычный issuer",
			issuerURL:    "http://localhost:8081/realms/datacomrade",
			wantJWKSURL:  "http://localhost:8081/realms/datacomrade/protocol/openid-connect/certs",
			wantCacheKey: "jwks:datacomrade",
		},
		{
			name:         "с завершающим слешем",
			issuerURL:    "http://localhost:8081/realms/datacomrade/",
			wantJWKSURL:  "http://localhost:8081/realms/datacomrade/protocol/openid-connect/certs",
			wantCacheKey: "jwks:datacomrade",
		},
		{
			name:         "с несколькими завершающими слешами",
			issuerURL:    "http://localhost:8081/realms/datacomrade///",
			wantJWKSURL:  "http://localhost:8081/realms/datacomrade/protocol/openid-connect/certs",
			wantCacheKey: "jwks:datacomrade",
		},
		{
			name:         "https и внешний хост",
			issuerURL:    "https://auth.example.com/realms/prod",
			wantJWKSURL:  "https://auth.example.com/realms/prod/protocol/openid-connect/certs",
			wantCacheKey: "jwks:prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewJWKSFetcher(nil, tt.issuerURL)

			if f.jwksURL != tt.wantJWKSURL {
				t.Errorf("jwksURL = %q, ожидался %q", f.jwksURL, tt.wantJWKSURL)
			}
			if f.cacheKey != tt.wantCacheKey {
				t.Errorf("cacheKey = %q, ожидался %q", f.cacheKey, tt.wantCacheKey)
			}
			if f.httpClient.Timeout == 0 {
				t.Error("httpClient без таймаута: запрос к Keycloak может повиснуть навсегда")
			}
		})
	}
}
