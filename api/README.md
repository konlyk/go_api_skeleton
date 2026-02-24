# OpenAPI Workflow

Source files:

- Root: `api/openapi.yaml`
- Paths: `api/paths/*.yaml`
- Schemas: `api/components/schemas/*.yaml`

Commands:

```bash
npm --prefix api install
npm --prefix api run lint
npm --prefix api run bundle
npm --prefix api run docs
```

Generate Go server/types into `api/openapi`:

```bash
go generate ./api/generate
```

Security model:

- Mark private operations with OpenAPI `security` requirements.
- Runtime validation/auth uses bundled contract (`api/openapi.bundled.yaml`) embedded by `api/spec.go`.
