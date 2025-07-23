#!/usr/bin/env just --justfile
run app="knox" *args="":
  go run ./cmd/{{app}} {{args}}

build:
  go build  -o ./bin ./cmd/...

test:
  go test -v ./internal/...

cover:
  go test ./...

bench:
    go test -benchmem -bench=. ./...

lint:
    golangci-lint run ./...

format:
    golangci-lint fmt ./...
