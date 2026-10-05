//go:build !darwin && !windows

package main

func registerGlobalHotkey(hk *Hotkey) bool {
	return false
}

func unregisterGlobalHotkey() {
}
