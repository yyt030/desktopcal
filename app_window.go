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

func (a *App) MinimiseWindow() {
	if a.mainWin != nil {
		a.mainWin.Minimise()
	}
}

func (a *App) MaximiseWindow() {
	if a.mainWin != nil {
		a.mainWin.ToggleMaximise()
	}
}

func (a *App) CloseWindow() {
	if a.mainWin != nil {
		a.mainWin.Close()
	}
}

func (a *App) GetVersion() string {
	return GetVersion()
}

func (a *App) GetWindowPosition() (int, int) {
	if a.mainWin != nil {
		return a.mainWin.Position()
	}
	return 0, 0
}

func (a *App) SetWindowPosition(x, y int) {
	if a.mainWin != nil {
		a.mainWin.SetPosition(x, y)
	}
}
