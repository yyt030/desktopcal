package main

import (
	"encoding/json"
	"fmt"
	"time"

	"desktopcal/internal/model"
)

type exportData struct {
	Events      []*model.Event    `json:"events"`
	Backgrounds map[string]string `json:"backgrounds"`
	ExportedAt  string            `json:"exportedAt"`
}

func (a *App) ExportData() (string, error) {
	events, err := a.db.GetAllEvents()
	if err != nil {
		return "", fmt.Errorf("导出事件失败: %w", err)
	}
	backgrounds, err := a.db.GetAllBackgrounds()
	if err != nil {
		return "", fmt.Errorf("导出背景失败: %w", err)
	}
	data := exportData{
		Events:      events,
		Backgrounds: backgrounds,
		ExportedAt:  time.Now().Format(time.RFC3339),
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (a *App) ImportData(jsonStr string) error {
	var data exportData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return fmt.Errorf("解析导入数据失败: %w", err)
	}
	if err := a.db.ImportEvents(data.Events); err != nil {
		return fmt.Errorf("导入事件失败: %w", err)
	}
	if err := a.db.ImportBackgrounds(data.Backgrounds); err != nil {
		return fmt.Errorf("导入背景失败: %w", err)
	}
	return nil
}

func (a *App) GetSettings() *model.Settings {
	s := &model.Settings{
		Opacity:        100,
		Width:          1100,
		Height:         750,
		Theme:          "dark",
		AlwaysTop:      true,
		FontSize:       13,
		FirstDayOfWeek: 0,
		ShowLunar:      true,
		ShowSolarTerm:  true,
		LayoutMode:     "compact",
	}
	keys := []string{"opacity", "width", "height", "theme", "always_top", "font_size", "first_day_of_week", "show_lunar", "show_solar_term", "layout_mode"}
	settings, err := a.db.GetSettings(keys)
	if err != nil {
		return s
	}
	if v, ok := settings["opacity"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &s.Opacity)
	}
	if v, ok := settings["width"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &s.Width)
	}
	if v, ok := settings["height"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &s.Height)
	}
	if v, ok := settings["theme"]; ok && v != "" {
		s.Theme = v
	}
	if v, ok := settings["always_top"]; ok && v != "" {
		s.AlwaysTop = v == "true"
	}
	if v, ok := settings["font_size"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &s.FontSize)
	}
	if v, ok := settings["first_day_of_week"]; ok && v != "" {
		fmt.Sscanf(v, "%d", &s.FirstDayOfWeek)
	}
	if v, ok := settings["show_lunar"]; ok && v != "" {
		s.ShowLunar = v == "true"
	}
	if v, ok := settings["show_solar_term"]; ok && v != "" {
		s.ShowSolarTerm = v == "true"
	}
	if v, ok := settings["layout_mode"]; ok && v != "" {
		s.LayoutMode = v
	}
	return s
}

func (a *App) SaveSettings(s *model.Settings) error {
	settings := map[string]string{
		"opacity":           fmt.Sprintf("%d", s.Opacity),
		"width":             fmt.Sprintf("%d", s.Width),
		"height":            fmt.Sprintf("%d", s.Height),
		"theme":             s.Theme,
		"font_size":         fmt.Sprintf("%d", s.FontSize),
		"first_day_of_week": fmt.Sprintf("%d", s.FirstDayOfWeek),
		"layout_mode":       s.LayoutMode,
		"always_top":        "false",
	}
	if s.AlwaysTop {
		settings["always_top"] = "true"
	}
	if s.ShowLunar {
		settings["show_lunar"] = "true"
	} else {
		settings["show_lunar"] = "false"
	}
	if s.ShowSolarTerm {
		settings["show_solar_term"] = "true"
	} else {
		settings["show_solar_term"] = "false"
	}
	return a.db.SetSettings(settings)
}
