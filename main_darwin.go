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
