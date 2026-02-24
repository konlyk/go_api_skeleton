package bootstrap

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func TestExtractBearerToken(t *testing.T) {
	t.Parallel()

	token, ok := extractBearerToken("Bearer abc123")
	if !ok {
		t.Fatal("expected bearer token to parse")
	}
	if token != "abc123" {
		t.Fatalf("expected token abc123, got %q", token)
	}

	if _, ok := extractBearerToken("Token abc123"); ok {
		t.Fatal("expected non-bearer token to fail")
	}
}

func TestOpenAPIAuthenticationFunc(t *testing.T) {
	t.Parallel()

	fn := newOpenAPIAuthenticationFunc("test-token")
	req := httptest.NewRequest(http.MethodGet, "/v1/hello", nil)
	req.Header.Set("Authorization", "Bearer test-token")

	err := fn(context.Background(), &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{Request: req},
		SecurityScheme: &openapi3.SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		},
	})
	if err != nil {
		t.Fatalf("expected auth to pass, got %v", err)
	}

	reqBad := httptest.NewRequest(http.MethodGet, "/v1/hello", nil)
	reqBad.Header.Set("Authorization", "Bearer wrong-token")

	err = fn(context.Background(), &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{Request: reqBad},
		SecurityScheme: &openapi3.SecurityScheme{
			Type:   "http",
			Scheme: "bearer",
		},
	})
	if err == nil {
		t.Fatal("expected auth to fail for wrong token")
	}
}
