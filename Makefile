APP_NAME = desktopcal
BUILD_DIR = build/bin
APP_DIR = $(BUILD_DIR)/$(APP_NAME).app
FRONTEND_DIR = frontend
# 编译期与链接期必须使用同一 macOS 最低版本，否则每个 C 对象都会触发
# "built for newer macOS version (12.0) than being linked (10.13)" 链接警告。
MACOS_MIN = 12.0
CGO_LDFLAGS = -mmacosx-version-min=$(MACOS_MIN) -framework UniformTypeIdentifiers -framework Security
CGO_CFLAGS = -mmacosx-version-min=$(MACOS_MIN)
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
