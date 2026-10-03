package main

// Window controls exposed to the frontend as bound service methods.

func (a *App) HideWindow() {
	if a.mainWin != nil {
		a.mainWin.Hide()
	}
}

func (a *App) ShowWindow() {
	if a.mainWin != nil {
		a.mainWin.Show()
	}
}

func (a *App) QuitApp() {
	if a.app != nil {
		a.app.Quit()
	}
}

func (a *App) SetWindowSize(width, height int) {
	if a.mainWin != nil {
		a.mainWin.SetSize(width, height)
	}
}

func (a *App) SetWindowAlwaysOnTop(onTop bool) {
	if a.mainWin != nil {
		a.mainWin.SetAlwaysOnTop(onTop)
	}
}
