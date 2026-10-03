package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"desktopcal/internal/db"
	"desktopcal/internal/lunar"
	"desktopcal/internal/model"
)

type App struct {
	ctx  context.Context
	db   *db.DB
	lcal *lunar.Calendar
}

func NewApp(database *db.DB) *App {
	return &App{
		db:   database,
		lcal: lunar.NewCalendar(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetMonthDays(year, month int) ([]*model.DayInfo, error) {
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	allEvents, err := a.db.GetEventsByMonth(year, month)
	if err != nil {
		allEvents = []*model.Event{}
	}
	eventsByDate := make(map[string][]*model.Event)
	for _, e := range allEvents {
		eventsByDate[e.Date] = append(eventsByDate[e.Date], e)
	}

	bgColors, err := a.db.GetDayBackgroundsByMonth(year, month)
	if err != nil {
		bgColors = make(map[string]string)
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	weekday := int(start.Weekday())
	gridStart := start.AddDate(0, 0, -weekday)

	days := make([]*model.DayInfo, 42)
	for i := 0; i < 42; i++ {
		d := gridStart.AddDate(0, 0, i)
		dateStr := d.Format("2006-01-02")

		ld := a.lcal.SolarToLunar(d)
		st := a.lcal.GetSolarTerm(d)

		events := eventsByDate[dateStr]
		if events == nil {
			events = []*model.Event{}
		}

		days[i] = &model.DayInfo{
			SolarDate: dateStr,
			Lunar:     ld,
			SolarTerm: st,
			Events:    events,
			IsToday:   dateStr == todayStr,
			IsCurrent: d.Month() == time.Month(month) && d.Year() == year,
			BgColor:   bgColors[dateStr],
		}
	}
	return days, nil
}

func (a *App) GetEventsByDate(date string) ([]*model.Event, error) {
	return a.db.GetEventsByDate(date)
}

func (a *App) CreateEvent(event *model.Event) (*model.Event, error) {
	if event.Title == "" {
		return nil, fmt.Errorf("事件标题不能为空")
	}
	if event.Date == "" {
		return nil, fmt.Errorf("事件日期不能为空")
	}
	return a.db.CreateEvent(event)
}

func (a *App) UpdateEvent(event *model.Event) error {
	return a.db.UpdateEvent(event)
}

func (a *App) DeleteEvent(id string) error {
	return a.db.DeleteEvent(id)
}

func (a *App) GetTodayLunar() *model.LunarDay {
	return a.lcal.SolarToLunar(time.Now())
}

func (a *App) GetTodaySolarTerm() string {
	return a.lcal.GetSolarTerm(time.Now())
}

func (a *App) SetDayBackground(date, color string) error {
	return a.db.SetDayBackground(date, color)
}

type exportData struct {
	Events       []*model.Event    `json:"events"`
	Backgrounds  map[string]string `json:"backgrounds"`
	ExportedAt   string            `json:"exportedAt"`
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
		Opacity:    100,
		Width:      1100,
		Height:     750,
		Theme:      "dark",
		AlwaysTop:  true,
		FontSize:   13,
	}
	if v, _ := a.db.GetSetting("opacity"); v != "" {
		fmt.Sscanf(v, "%d", &s.Opacity)
	}
	if v, _ := a.db.GetSetting("width"); v != "" {
		fmt.Sscanf(v, "%d", &s.Width)
	}
	if v, _ := a.db.GetSetting("height"); v != "" {
		fmt.Sscanf(v, "%d", &s.Height)
	}
	if v, _ := a.db.GetSetting("theme"); v != "" {
		s.Theme = v
	}
	if v, _ := a.db.GetSetting("always_top"); v != "" {
		s.AlwaysTop = v == "true"
	}
	if v, _ := a.db.GetSetting("font_size"); v != "" {
		fmt.Sscanf(v, "%d", &s.FontSize)
	}
	return s
}

func (a *App) SaveSettings(s *model.Settings) error {
	a.db.SetSetting("opacity", fmt.Sprintf("%d", s.Opacity))
	a.db.SetSetting("width", fmt.Sprintf("%d", s.Width))
	a.db.SetSetting("height", fmt.Sprintf("%d", s.Height))
	a.db.SetSetting("theme", s.Theme)
	a.db.SetSetting("font_size", fmt.Sprintf("%d", s.FontSize))
	if s.AlwaysTop {
		a.db.SetSetting("always_top", "true")
	} else {
		a.db.SetSetting("always_top", "false")
	}
	return nil
}
