.PHONY: openapi-install openapi-lint openapi-bundle openapi-docs openapi-generate generate run test

openapi-install:
	npm --prefix api install

openapi-lint:
	npm --prefix api run lint

openapi-bundle:
	npm --prefix api run bundle

openapi-docs:
	npm --prefix api run docs

openapi-generate:
	go generate ./api/generate

generate: openapi-bundle openapi-generate

run:
	go run ./cmd

test:
	go test ./...
