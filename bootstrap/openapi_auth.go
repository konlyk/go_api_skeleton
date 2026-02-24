package bootstrap

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func newOpenAPIAuthenticationFunc(expectedToken string) openapi3filter.AuthenticationFunc {
	return func(_ context.Context, input *openapi3filter.AuthenticationInput) error {
		if input == nil || input.RequestValidationInput == nil || input.RequestValidationInput.Request == nil {
			return errors.New("missing auth input")
		}
		if input.SecurityScheme == nil {
			return input.NewError(errors.New("missing security scheme"))
		}

		if !strings.EqualFold(input.SecurityScheme.Type, "http") || !strings.EqualFold(input.SecurityScheme.Scheme, "bearer") {
			return input.NewError(errors.New("unsupported security scheme"))
		}

		token, ok := extractBearerToken(input.RequestValidationInput.Request.Header.Get("Authorization"))
		if !ok {
			return input.NewError(errors.New("missing or invalid Authorization header"))
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) != 1 {
			return input.NewError(errors.New("invalid bearer token"))
		}

		return nil
	}
}

func extractBearerToken(header string) (string, bool) {
	parts := strings.Fields(strings.TrimSpace(header))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}
