package main

import (
	"embed"
	"log"
	"sync/atomic"

	"desktopcal/internal/db"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

// quitting 标记应用正在退出：窗口关闭 hook 据此放行，
// 避免 Quit 的清理流程被「只隐藏窗口」的拦截逻辑卡住。
var quitting atomic.Bool

func main() {
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	service := NewApp(database)

	app := application.New(application.Options{
		Name:     "桌面日历",
		Services: []application.Service{application.NewService(service)},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "桌面日历",
		Width:            1100,
		Height:           750,
		MinWidth:         900,
		MinHeight:        600,
		Hidden:           false,
		Frameless:        true,
		BackgroundColour: application.RGBA{Red: 30, Green: 30, Blue: 46, Alpha: 255},
		KeyBindings: map[string]func(application.Window){
			"cmd+w": func(w application.Window) { w.Hide() },
		},
	})
	service.setRuntimeRefs(app, win)

	// 关闭按钮只隐藏窗口，应用保留在系统托盘中。
	win.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		if quitting.Load() {
			return
		}
		e.Cancel()
		win.Hide()
	})

	setupSystemTray(app, win)
	setupPlatformEvents(app, win)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// setupSystemTray 创建菜单栏图标：左键切换主窗口显隐，右键弹出菜单。
// 有意不使用 AttachWindow——那会把窗口强制挪到图标正下方并提升到弹层
// 窗口级别，破坏用户对窗口位置的摆放。
func setupSystemTray(app *application.App, win *application.WebviewWindow) {
	menu := app.NewMenu()
	menu.Add("显示日历").OnClick(func(*application.Context) { win.Show(); win.Focus() })
	menu.Add("隐藏窗口").OnClick(func(*application.Context) { win.Hide() })
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(*application.Context) {
		quitting.Store(true)
		app.Quit()
	})

	tray := app.SystemTray.New()
	setupPlatformTray(tray)
	tray.SetMenu(menu)
	tray.OnClick(func() {
		// 可见且聚焦 → 隐藏；否则唤到前台（被其他窗口遮挡时也能提到最前）。
		if win.IsVisible() && win.IsFocused() {
			win.Hide()
		} else {
			win.Show()
			win.Focus()
		}
	})
}
