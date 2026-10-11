package main

import (
	"desktopcal/internal/db"
	"desktopcal/internal/lunar"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// appInstance 全局引用，供快捷键回调等外部函数访问。
var appInstance *App

// App 是暴露给前端的绑定服务，聚合数据库、农历计算与窗口控制能力。
type App struct {
	app     *application.App
	mainWin application.Window
	db      *db.DB
	lcal    *lunar.Calendar
}

func NewApp(database *db.DB) *App {
	a := &App{
		db:   database,
		lcal: lunar.NewCalendar(),
	}
	appInstance = a
	return a
}

// setRuntimeRefs is called from main once the Wails app and main window exist,
// giving the bound service access to window controls. Unexported so it is not
// generated as a frontend-bound method.
func (a *App) setRuntimeRefs(app *application.App, win application.Window) {
	a.app = app
	a.mainWin = win
}

// GetPlatform 返回当前运行平台（darwin/windows/linux），供前端判断功能可用性。
func (a *App) GetPlatform() string {
	return runtime.GOOS
}
