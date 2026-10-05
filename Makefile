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
# go-winres：生成 Windows 资源文件（图标、清单）
WINRES ?= $(HOME)/.gvm/pkgsets/go1.25.14/global/bin/go-winres

# 版本信息：从 git 获取，编译时注入
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS = -w -s -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)

.PHONY: dev build frontend clean run bundle bindings windows winres

run: bundle
	open $(APP_DIR)

bundle: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" GOMAXPROCS=2 nice -n 15 $(GO) build -p 2 -tags production -ldflags="$(LDFLAGS)" -o $(APP_DIR)/Contents/MacOS/$(APP_NAME) .

frontend:
	test -d $(FRONTEND_DIR)/node_modules || (cd $(FRONTEND_DIR) && $(NPM) install)
	cd $(FRONTEND_DIR) && nice -n 15 $(NPX) vite build
	mkdir -p $(APP_DIR)/Contents/MacOS
	mkdir -p $(APP_DIR)/Contents/Resources
	cp build/darwin/Info.plist $(APP_DIR)/Contents/Info.plist
	cp build/darwin/icon.icns $(APP_DIR)/Contents/Resources/icon.icns

build: frontend
	CGO_ENABLED=1 CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" GOROOT="$(HOME)/.gvm/gos/go1.25.14" GOMAXPROCS=2 nice -n 15 $(GO) build -p 2 -tags production -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) .

# Windows 交叉编译：-H windowsgui 隐藏 CMD 黑框
windows: frontend winres
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 GOROOT="$(HOME)/.gvm/gos/go1.25.14" GOMAXPROCS=2 nice -n 15 $(GO) build -p 2 -tags production -ldflags="$(LDFLAGS) -H windowsgui" -o $(BUILD_DIR)/$(APP_NAME).exe .

# 生成 Windows 资源文件（图标、清单）。修改 winres/winres.json 后需重新执行。
winres:
	$(WINRES) make

clean:
	rm -rf $(BUILD_DIR) $(FRONTEND_DIR)/dist

# 重新生成前端绑定。改动 internal/model 或绑定方法后需手动执行。
# 必须带 -ts，否则默认输出带 JSDoc 的 .js。不挂到 build 依赖上，
# 以免在没有 wails3 的环境里构建失败。
bindings:
	PATH="$(HOME)/.gvm/gos/go1.25.14/bin:$$PATH" GOROOT="$(HOME)/.gvm/gos/go1.25.14" \
		$(WAILS3) generate bindings -ts ./...
