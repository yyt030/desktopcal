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
# Wails v3 CLI。注意与 v2 的 `wails` 不是同一个程序，本工程只能用 wails3。
WAILS3 ?= $(HOME)/.gvm/pkgsets/go1.25.14/global/bin/wails3

.PHONY: dev build frontend clean run bundle bindings

run: bundle
	open $(APP_DIR)

bundle: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" GOMAXPROCS=2 $(GO) build -p 2 -tags production -ldflags="-w -s" -o $(APP_DIR)/Contents/MacOS/$(APP_NAME) .

frontend:
	test -d $(FRONTEND_DIR)/node_modules || (cd $(FRONTEND_DIR) && $(NPM) install)
	cd $(FRONTEND_DIR) && $(NPX) vite build
	mkdir -p $(APP_DIR)/Contents/MacOS
	mkdir -p $(APP_DIR)/Contents/Resources
	cp build/darwin/Info.plist $(APP_DIR)/Contents/Info.plist
	cp build/darwin/icon.icns $(APP_DIR)/Contents/Resources/icon.icns

build: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" GOMAXPROCS=2 $(GO) build -p 2 -tags production -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

clean:
	rm -rf $(BUILD_DIR) $(FRONTEND_DIR)/dist

# 重新生成前端绑定。改动 internal/model 或绑定方法后需手动执行。
# 必须带 -ts，否则默认输出带 JSDoc 的 .js。不挂到 build 依赖上，
# 以免在没有 wails3 的环境里构建失败。
bindings:
	PATH="$(HOME)/.gvm/gos/go1.25.14/bin:$$PATH" GOROOT="$(HOME)/.gvm/gos/go1.25.14" \
		$(WAILS3) generate bindings -ts ./...
