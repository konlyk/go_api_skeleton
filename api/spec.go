package api

import (
	_ "embed"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	//go:embed openapi.bundled.yaml
	bundledOpenAPI []byte
)

func LoadBundledSpec() (*openapi3.T, error) {
	loader := openapi3.NewLoader()

	spec, err := loader.LoadFromData(bundledOpenAPI)
	if err != nil {
		return nil, fmt.Errorf("load bundled openapi spec: %w", err)
	}

	return spec, nil
}
