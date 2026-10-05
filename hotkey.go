package main

import (
	"strings"
)

// Hotkey 表示一个全局快捷键配置。
type Hotkey struct {
	Key        string // 主键，如 "c"
	Cmd        bool   // macOS: Cmd 键
	Ctrl       bool   // Ctrl 键
	Alt        bool   // Alt/Option 键
	Shift      bool   // Shift 键
}

// ParseHotkey 解析快捷键字符串，格式为 "mod1+mod2+key"。
// 支持的修饰键：cmd, ctrl, alt, shift（不区分大小写）。
func ParseHotkey(s string) *Hotkey {
	if s == "" {
		return nil
	}
	parts := strings.Split(strings.ToLower(s), "+")
	if len(parts) < 2 {
		return nil
	}
	hk := &Hotkey{}
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if i == len(parts)-1 {
			hk.Key = part
		} else {
			switch part {
			case "cmd", "command":
				hk.Cmd = true
			case "ctrl", "control":
				hk.Ctrl = true
			case "alt", "option":
				hk.Alt = true
			case "shift":
				hk.Shift = true
			}
		}
	}
	if hk.Key == "" {
		return nil
	}
	return hk
}

// String 返回快捷键的字符串表示。
func (h *Hotkey) String() string {
	if h == nil {
		return ""
	}
	var parts []string
	if h.Cmd {
		parts = append(parts, "cmd")
	}
	if h.Ctrl {
		parts = append(parts, "ctrl")
	}
	if h.Alt {
		parts = append(parts, "alt")
	}
	if h.Shift {
		parts = append(parts, "shift")
	}
	parts = append(parts, h.Key)
	return strings.Join(parts, "+")
}
