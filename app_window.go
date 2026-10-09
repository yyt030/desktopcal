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

func (a *App) GetDatabasePath() string {
	if a.db != nil {
		return a.db.Path()
	}
	return ""
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

func (a *App) GetWindowSize() (int, int) {
	if a.mainWin != nil {
		return a.mainWin.Size()
	}
	return 0, 0
}

func (a *App) GetScreenWorkArea() (x, y, width, height int) {
	if a.mainWin == nil {
		return 0, 0, 1920, 1080
	}
	screen, err := a.mainWin.GetScreen()
	if err != nil || screen == nil {
		return 0, 0, 1920, 1080
	}
	return screen.WorkArea.X, screen.WorkArea.Y, screen.WorkArea.Width, screen.WorkArea.Height
}
