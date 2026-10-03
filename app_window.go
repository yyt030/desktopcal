package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}
