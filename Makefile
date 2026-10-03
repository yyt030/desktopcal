APP_NAME = desktopcal
BUILD_DIR = build/bin
FRONTEND_DIR = frontend
CGO_LDFLAGS = -framework UniformTypeIdentifiers
GO ?= $(HOME)/.gvm/gos/go1.25.14/bin/go
NPM ?= npm
NPX ?= npx

.PHONY: dev build frontend clean run

run: frontend
	CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" $(GO) build -tags dev -o $(BUILD_DIR)/$(APP_NAME) .
	$(BUILD_DIR)/$(APP_NAME)

build: frontend
	CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" $(GO) build -tags production -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

frontend:
	cd $(FRONTEND_DIR) && $(NPM) install && $(NPX) vite build

clean:
	rm -rf $(BUILD_DIR) $(FRONTEND_DIR)/dist
