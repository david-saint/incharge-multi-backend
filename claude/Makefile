.PHONY: build run dev test clean lint fmt tidy docker-up docker-down docker-build

# ── Variables ────────────────────────────────────────────────────────────────
BINARY  := incharge
CMD     := ./cmd/server
GOFLAGS := -ldflags="-s -w"

# ── Build ────────────────────────────────────────────────────────────────────
build:
	go build $(GOFLAGS) -o $(BINARY) $(CMD)

run: build
	./$(BINARY)

dev:
	@echo "Starting with live reload (requires air: go install github.com/air-verse/air@latest)"
	air -c .air.toml 2>/dev/null || go run $(CMD)

# ── Code Quality ─────────────────────────────────────────────────────────────
fmt:
	gofmt -w -s .

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

vet:
	go vet ./...

# ── Testing ──────────────────────────────────────────────────────────────────
test:
	go test -v -race -count=1 ./...

test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ── Docker ───────────────────────────────────────────────────────────────────
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

# ── Cleanup ──────────────────────────────────────────────────────────────────
clean:
	rm -f $(BINARY) server coverage.out coverage.html
	go clean

# ── Help ─────────────────────────────────────────────────────────────────────
help:
	@echo "Usage:"
	@echo "  make build         Build the server binary"
	@echo "  make run           Build and run the server"
	@echo "  make dev           Run with live reload (air) or fallback to go run"
	@echo "  make fmt           Format Go source files"
	@echo "  make lint          Run golangci-lint"
	@echo "  make tidy          Run go mod tidy"
	@echo "  make vet           Run go vet"
	@echo "  make test          Run all tests"
	@echo "  make test-cover    Run tests with coverage report"
	@echo "  make docker-build  Build Docker image"
	@echo "  make docker-up     Start containers (db + app)"
	@echo "  make docker-down   Stop containers"
	@echo "  make docker-logs   Tail app container logs"
	@echo "  make clean         Remove build artifacts"
