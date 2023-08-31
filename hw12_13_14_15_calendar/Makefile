BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run --out-format=github-actions ./...

.PHONY: models
models:
	@sqlboiler --add-soft-deletes -c sqlboiler.toml psql
	@echo "Models generated successfully"

.PHONY: postgres_migration_up
postgres_migration_up:
	@migrate -path database/migrations -database postgresql://postgres:postgres@postgres:5432/otus?sslmode=disable\&x-migrations-table=calendar_migrations up

postgres_migration_create: NAME=$(call getArgs)
postgres_migration_create:
	migrate create -ext sql -dir database/migrations/ ${NAME}

.PHONY: build run build-img run-img version test lint

# Функция получения аргументов
define getArgs
$(strip $(filter-out $@,$(MAKECMDGOALS)))
endef