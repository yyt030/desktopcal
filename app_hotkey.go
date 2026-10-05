package main

import "desktopcal/internal/model"

// GetHotkey 返回当前全局快捷键配置。
func (a *App) GetHotkey() string {
	settings := a.GetSettings()
	return settings.Hotkey
}

// SetHotkey 更新全局快捷键并重新注册。
func (a *App) SetHotkey(hotkey string) error {
	settings := a.GetSettings()
	settings.Hotkey = hotkey
	if err := a.SaveSettings(settings); err != nil {
		return err
	}
	// 重新注册快捷键
	unregisterGlobalHotkey()
	if hk := ParseHotkey(hotkey); hk != nil {
		registerGlobalHotkey(hk)
	}
	return nil
}

// ResetHotkey 重置为默认快捷键。
func (a *App) ResetHotkey() error {
	return a.SetHotkey(model.DefaultHotkey())
}
