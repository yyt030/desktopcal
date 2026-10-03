APP_NAME = desktopcal
BUILD_DIR = build/bin
APP_DIR = $(BUILD_DIR)/$(APP_NAME).app
FRONTEND_DIR = frontend
CGO_LDFLAGS = -framework UniformTypeIdentifiers -framework Security
CGO_CFLAGS = -mmacosx-version-min=12.0
GO ?= $(HOME)/.gvm/gos/go1.25.14/bin/go
NPM ?= npm
NPX ?= npx

.PHONY: dev build frontend clean run bundle

run: bundle
	open $(APP_DIR)

bundle: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" $(GO) build -tags production -ldflags="-w -s" -o $(APP_DIR)/Contents/MacOS/$(APP_NAME) .

frontend:
	cd $(FRONTEND_DIR) && $(NPM) install && $(NPX) vite build
	mkdir -p $(APP_DIR)/Contents/MacOS
	mkdir -p $(APP_DIR)/Contents/Resources
	cp build/darwin/Info.plist $(APP_DIR)/Contents/Info.plist
	cp build/darwin/icon.icns $(APP_DIR)/Contents/Resources/icon.icns

build: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" $(GO) build -tags production -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

clean:
	rm -rf $(BUILD_DIR) $(FRONTEND_DIR)/dist
