package main

import (
	"fmt"
	"strconv"

	"desktopcal/internal/model"
)

// calendarSettingKeys 是渲染月视图所需的最小设置集合。
var calendarSettingKeys = []string{"layout_mode", "first_day_of_week"}

// settingCodec 描述一个设置项在 SQLite 字符串值与 model.Settings 字段之间的转换。
// 新增设置项时只需在 settingsSpec 中追加一行，读写两侧会自动同步。
type settingCodec interface {
	key() string
	encode(dst *model.Settings) string
	decode(dst *model.Settings, raw string)
}

type fieldCodec[T any] struct {
	name   string
	field  func(*model.Settings) *T
	format func(T) string
	parse  func(raw string, fallback T) T
}

func (c fieldCodec[T]) key() string { return c.name }

func (c fieldCodec[T]) encode(dst *model.Settings) string {
	return c.format(*c.field(dst))
}

// decode 在值缺失或非法时保留字段当前的默认值。
func (c fieldCodec[T]) decode(dst *model.Settings, raw string) {
	*c.field(dst) = c.parse(raw, *c.field(dst))
}

func intSetting(name string, field func(*model.Settings) *int) settingCodec {
	return fieldCodec[int]{
		name:   name,
		field:  field,
		format: strconv.Itoa,
		parse: func(raw string, fallback int) int {
			if v, err := strconv.Atoi(raw); err == nil {
				return v
			}
			return fallback
		},
	}
}

func boolSetting(name string, field func(*model.Settings) *bool) settingCodec {
	return fieldCodec[bool]{
		name:   name,
		field:  field,
		format: strconv.FormatBool,
		parse: func(raw string, fallback bool) bool {
			if v, err := strconv.ParseBool(raw); err == nil {
				return v
			}
			return fallback
		},
	}
}

func stringSetting(name string, field func(*model.Settings) *string) settingCodec {
	return fieldCodec[string]{
		name:   name,
		field:  field,
		format: func(v string) string { return v },
		parse: func(raw string, fallback string) string {
			if raw == "" {
				return fallback
			}
			return raw
		},
	}
}

// settingsSpec 是所有可持久化设置项的唯一声明，读写共用。
var settingsSpec = []settingCodec{
	intSetting("opacity", func(s *model.Settings) *int { return &s.Opacity }),
	intSetting("width", func(s *model.Settings) *int { return &s.Width }),
	intSetting("height", func(s *model.Settings) *int { return &s.Height }),
	intSetting("font_size", func(s *model.Settings) *int { return &s.FontSize }),
	stringSetting("font_family", func(s *model.Settings) *string { return &s.FontFamily }),
	intSetting("calendar_font_size", func(s *model.Settings) *int { return &s.CalendarFontSize }),
	stringSetting("calendar_font_family", func(s *model.Settings) *string { return &s.CalendarFontFamily }),
	intSetting("first_day_of_week", func(s *model.Settings) *int { return &s.FirstDayOfWeek }),
	stringSetting("theme", func(s *model.Settings) *string { return &s.Theme }),
	stringSetting("layout_mode", func(s *model.Settings) *string { return &s.LayoutMode }),
	stringSetting("hotkey", func(s *model.Settings) *string { return &s.Hotkey }),
	stringSetting("export_path", func(s *model.Settings) *string { return &s.ExportPath }),
	boolSetting("always_top", func(s *model.Settings) *bool { return &s.AlwaysTop }),
	boolSetting("show_lunar", func(s *model.Settings) *bool { return &s.ShowLunar }),
	boolSetting("show_solar_term", func(s *model.Settings) *bool { return &s.ShowSolarTerm }),
	boolSetting("close_to_tray", func(s *model.Settings) *bool { return &s.CloseToTray }),
	boolSetting("auto_start", func(s *model.Settings) *bool { return &s.AutoStart }),
	intSetting("editor_width", func(s *model.Settings) *int { return &s.EditorWidth }),
	intSetting("editor_height", func(s *model.Settings) *int { return &s.EditorHeight }),
}

func (a *App) GetSettings() *model.Settings {
	s := model.DefaultSettings()
	stored, err := a.db.GetSettings(settingKeys())
	if err != nil {
		return s
	}
	for _, item := range settingsSpec {
		if raw, ok := stored[item.key()]; ok && raw != "" {
			item.decode(s, raw)
		}
	}
	return s
}

func (a *App) SaveSettings(s *model.Settings) error {
	if s == nil {
		return fmt.Errorf("设置不能为空")
	}
	old := a.GetSettings()
	values := make(map[string]string, len(settingsSpec))
	for _, item := range settingsSpec {
		values[item.key()] = item.encode(s)
	}
	if err := a.db.SetSettings(values); err != nil {
		return err
	}
	if old.Hotkey != s.Hotkey {
		unregisterGlobalHotkey()
		if hk := ParseHotkey(s.Hotkey); hk != nil {
			registerGlobalHotkey(hk)
		}
	}
	if old.AutoStart != s.AutoStart {
		setAutoStart(s.AutoStart)
	}
	return nil
}

// settingKeys 返回全部设置项的键名，用于一次性批量读取。
func settingKeys() []string {
	keys := make([]string, len(settingsSpec))
	for i, item := range settingsSpec {
		keys[i] = item.key()
	}
	return keys
}

func (a *App) SaveEditorSize(width, height int) error {
	if width < 300 || height < 200 {
		return nil
	}
	s := a.GetSettings()
	s.EditorWidth = width
	s.EditorHeight = height
	return a.SaveSettings(s)
}

// ResetSettings 清空全部持久化设置并返回默认值，同时触发快捷键/自启动等副作用。
func (a *App) ResetSettings() (*model.Settings, error) {
	old := a.GetSettings()
	if err := a.db.ClearSettings(); err != nil {
		return nil, err
	}
	def := model.DefaultSettings()
	if old.Hotkey != def.Hotkey {
		unregisterGlobalHotkey()
		if hk := ParseHotkey(def.Hotkey); hk != nil {
			registerGlobalHotkey(hk)
		}
	}
	if old.AutoStart != def.AutoStart {
		setAutoStart(def.AutoStart)
	}
	return def, nil
}
