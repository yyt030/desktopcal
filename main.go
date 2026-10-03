package main

import (
	"embed"
	"log"

	"desktopcal/internal/db"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/darwin/appicon.png
var appIcon []byte

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
		BackgroundColour: application.RGBA{Red: 30, Green: 30, Blue: 46, Alpha: 255},
		KeyBindings: map[string]func(application.Window){
			"cmd+w": func(w application.Window) { w.Hide() },
		},
	})
	service.setRuntimeRefs(app, win)

	// 关闭按钮只隐藏窗口，应用保留在系统托盘中。
	win.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		e.Cancel()
		win.Hide()
	})

	// 原生系统托盘：点击显示窗口，右键菜单显示/退出。
	menu := app.NewMenu()
	menu.Add("显示日历").OnClick(func(*application.Context) { win.Show(); win.Focus() })
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(*application.Context) { app.Quit() })

	_ = menu
	// TEMP: tray disabled for isolation
	// tray := app.SystemTray.New()
	// tray.SetIcon(appIcon)
	// tray.SetTooltip("桌面日历")
	// tray.SetMenu(menu)
	// tray.AttachWindow(win).WindowOffset(5)

	// Dock 图标重新点击时恢复窗口。
	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(e *application.ApplicationEvent) {
		win.Show()
		win.Focus()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
