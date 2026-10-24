# Команды для сборки и проверки backend.
GO ?= go
GOLANGCI_LINT ?= golangci-lint
LINTER_VERSION ?= v2.14.0
ARGS ?=
BINARY ?= bin/smart-warehouse

.DEFAULT_GOAL := help
.PHONY: help install build run format lint test

help: ## Показать список команд
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Скачать зависимости приложения и установить линтер
	$(GO) mod download
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINTER_VERSION)

build: ## Собрать backend
	@mkdir -p $(dir $(BINARY))
	$(GO) build -trimpath -o $(BINARY) ./cmd/app

run: ## Запустить backend с переменными окружения
	$(GO) run ./cmd/app $(ARGS)

format: ## Применить gofumpt и goimports
	$(GOLANGCI_LINT) fmt

lint: ## Проверить зависимости, форматирование и код
	$(GO) mod tidy -diff
	$(GOLANGCI_LINT) config verify
	$(GOLANGCI_LINT) fmt --diff
	$(GOLANGCI_LINT) run

test: ## Запустить тесты с race detector и покрытием
	$(GO) test -race -coverpkg=./internal/... -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
