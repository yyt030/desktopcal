APP_NAME = desktopcal
BUILD_DIR = build/bin
FRONTEND_DIR = frontend
CGO_LDFLAGS = -framework UniformTypeIdentifiers

.PHONY: dev build frontend clean run

run: frontend
	CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -tags dev -o $(BUILD_DIR)/$(APP_NAME) .
	$(BUILD_DIR)/$(APP_NAME)

build: frontend
	CGO_ENABLED=1 CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -tags production -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) .

frontend:
	cd $(FRONTEND_DIR) && npm install && npx vite build

clean:
	rm -rf $(BUILD_DIR) $(FRONTEND_DIR)/dist
