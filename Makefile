APP_NAME    := grapevine
GO          := go
GOFLAGS     := -ldflags="-s -w"
WEB_DIR     := web
DIST_DIR    := $(WEB_DIR)/dist

.PHONY: all build run clean dev \
        build-backend build-frontend \
        lint test tidy \
        docker docker-up docker-down

# ─── Default ──────────────────────────────────────────────
all: build

# ─── Build ────────────────────────────────────────────────
build: build-frontend build-backend

build-backend:
	$(GO) build $(GOFLAGS) -o $(APP_NAME) ./cmd/grapevine

build-frontend:
	cd $(WEB_DIR) && npm install && npm run build

# ─── Run ──────────────────────────────────────────────────
run: build-backend
	./$(APP_NAME)

dev:
	$(GO) run ./cmd/grapevine

# ─── Quality ──────────────────────────────────────────────
tidy:
	$(GO) mod tidy

lint:
	$(GO) vet ./...

test:
	$(GO) test -v -race ./...

# ─── Docker ───────────────────────────────────────────────
docker:
	docker build -t $(APP_NAME) .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# ─── Clean ────────────────────────────────────────────────
clean:
	rm -f $(APP_NAME)
	rm -rf $(DIST_DIR)
