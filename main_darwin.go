//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setDockIconVisible(bool visible) {
    NSApplication *app = [NSApplication sharedApplication];
    if (visible) {
        [app setActivationPolicy:NSApplicationActivationPolicyRegular];
    } else {
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];
    }
}
*/
import "C"

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed build/darwin/trayicon.png
var trayIcon []byte

func setupPlatformTray(tray *application.SystemTray) {
	// macOS 模板图标（纯黑 + alpha）：系统按菜单栏明暗自动反色
	tray.SetTemplateIcon(trayIcon)
}

func setupPlatformEvents(app *application.App, win *application.WebviewWindow) {
	// Dock 图标重新点击时恢复窗口
	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(e *application.ApplicationEvent) {
		win.Show()
		win.Focus()
	})
}

func setDockIconVisible(visible bool) {
	C.setDockIconVisible(C.bool(visible))
}

func shouldStartHidden() bool {
	return false
}

// setAutoStart 通过 LaunchAgent 管理 macOS 开机自启动。
// 在 ~/Library/LaunchAgents/ 下创建或移除 plist 文件。
func setAutoStart(enabled bool) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	plistDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	plistPath := filepath.Join(plistDir, "com.desktopcal.plist")

	if !enabled {
		os.Remove(plistPath)
		return nil
	}

	if err := os.MkdirAll(plistDir, 0755); err != nil {
		return err
	}

	// 获取当前可执行文件路径（.app bundle 内的二进制）
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// macOS 下使用 open -a 打开 app bundle
	// exe 通常是 /path/to/App.app/Contents/MacOS/app
	// 需要找到 .app 目录
	appPath := exe
	for i := len(exe) - 1; i >= 0; i-- {
		if exe[i] == '/' {
			if i >= 4 && exe[i-4:i] == ".app" {
				appPath = exe[:i+1]
				break
			}
		}
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.desktopcal</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/bin/open</string>
        <string>-a</string>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
`, appPath)

	return os.WriteFile(plistPath, []byte(plist), 0644)
}
