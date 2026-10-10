//go:build windows

package main

import (
	_ "embed"
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/sys/windows"
)

//go:embed build/windows/trayicon.png
var trayIcon []byte

const (
	wsExToolWindow = 0x00000080
	wsExAppWindow  = 0x00040000
)

var (
	user32DLL           = windows.NewLazySystemDLL("user32.dll")
	getWindowLongProc   = user32DLL.NewProc("GetWindowLongPtrW")
	setWindowLongProc   = user32DLL.NewProc("SetWindowLongPtrW")
	enumWindowsProc     = user32DLL.NewProc("EnumWindows")
	getWindowTextProc   = user32DLL.NewProc("GetWindowTextW")
	getTextLengthProc   = user32DLL.NewProc("GetWindowTextLengthW")
	getWindowThreadProc = user32DLL.NewProc("GetWindowThreadProcessId")
)

var targetHwnd windows.Handle
var targetProcessID uint32

func enumWindowsCallback(hwnd windows.Handle, lParam uintptr) uintptr {
	var processID uint32
	getWindowThreadProc.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processID)))
	
	if processID == targetProcessID {
		// 检查窗口标题
		length, _, _ := getTextLengthProc.Call(uintptr(hwnd))
		if length > 0 {
			buf := make([]uint16, length+1)
			getWindowTextProc.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), length+1)
			title := windows.UTF16ToString(buf)
			if title == "桌面日历" {
				targetHwnd = hwnd
				return 0 // 停止枚举
			}
		}
	}
	return 1 // 继续枚举
}

func setupPlatformTray(tray *application.SystemTray) {
	// Windows 使用普通图标（非模板）
	tray.SetIcon(trayIcon)
	tray.SetTooltip("桌面日历")
}

func setupPlatformEvents(app *application.App, win *application.WebviewWindow) {
	// Windows 无 Dock 事件，无需处理
}

func shouldStartHidden() bool {
	return false
}

func setDockIconVisible(visible bool) {
	// 获取当前进程 ID
	currentProcessID := uint32(windows.GetCurrentProcessId())
	targetProcessID = currentProcessID
	targetHwnd = 0
	
	// 枚举窗口找到主窗口
	enumWindowsProc.Call(
		windows.NewCallback(enumWindowsCallback),
		0,
	)
	
	if targetHwnd == 0 {
		return
	}

	// GWL_EXSTYLE = -20
	const gwlExStyle = ^uintptr(19)
	
	// 获取当前扩展样式
	style, _, _ := getWindowLongProc.Call(uintptr(targetHwnd), gwlExStyle)
	
	if visible {
		// 恢复为普通窗口（显示在任务栏）
		style |= wsExAppWindow
		style &^= wsExToolWindow
	} else {
		// 设置为工具窗口（不显示在任务栏）
		style |= wsExToolWindow
		style &^= wsExAppWindow
	}
	
	setWindowLongProc.Call(uintptr(targetHwnd), gwlExStyle, style)
}
