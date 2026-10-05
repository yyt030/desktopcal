//go:build windows

package main

import (
	_ "embed"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed build/windows/trayicon.png
var trayIcon []byte

func setupPlatformTray(tray *application.SystemTray) {
	// Windows 使用普通图标（非模板）
	tray.SetIcon(trayIcon)
	tray.SetTooltip("桌面日历")
}

func setupPlatformEvents(app *application.App, win *application.WebviewWindow) {
	// Windows 无 Dock 事件，无需处理
}
