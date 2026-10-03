package main

import (
	"fmt"
	"time"

	"desktopcal/internal/db"
	"desktopcal/internal/lunar"
	"desktopcal/internal/model"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	app     *application.App
	mainWin application.Window
	db      *db.DB
	lcal    *lunar.Calendar
}

func NewApp(database *db.DB) *App {
	return &App{
		db:   database,
		lcal: lunar.NewCalendar(),
	}
}

// setRuntimeRefs is called from main once the Wails app and main window exist,
// giving the bound service access to window controls. Unexported so it is not
// generated as a frontend-bound method.
func (a *App) setRuntimeRefs(app *application.App, win application.Window) {
	a.app = app
	a.mainWin = win
}

func (a *App) GetMonthDays(year, month int) ([]*model.DayInfo, error) {
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	firstDayOfWeek := 0
	if v, _ := a.db.GetSetting("first_day_of_week"); v != "" {
		fmt.Sscanf(v, "%d", &firstDayOfWeek)
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)
	weekday := int(start.Weekday())
	offset := (weekday - firstDayOfWeek + 7) % 7
	naturalStart := start.AddDate(0, 0, -offset)
	lastDayOffset := (int(end.Weekday()) - firstDayOfWeek + 7) % 7
	trailing := (6 - lastDayOffset) % 7
	totalDays := offset + end.Day() + trailing

	todayRow := 0
	if today.Year() == year && int(today.Month()) == month {
		todayRow = (offset + today.Day() - 1) / 7
	}
	rows := totalDays / 7
	leading := rows/2 - todayRow
	if leading < 0 {
		leading = 0
	}
	if leading > rows {
		leading = rows
	}
	gridStart := naturalStart.AddDate(0, 0, -leading*7)
	gridEnd := gridStart.AddDate(0, 0, totalDays-1)

	allEvents, err := a.db.GetEventsByDateRange(gridStart.Format("2006-01-02"), gridEnd.Format("2006-01-02"))
	if err != nil {
		allEvents = []*model.Event{}
	}
	eventsByDate := make(map[string][]*model.Event)
	for _, e := range allEvents {
		eventsByDate[e.Date] = append(eventsByDate[e.Date], e)
	}

	bgColors, err := a.db.GetDayBackgroundsByDateRange(gridStart.Format("2006-01-02"), gridEnd.Format("2006-01-02"))
	if err != nil {
		bgColors = make(map[string]string)
	}

	days := make([]*model.DayInfo, totalDays)
	for i := 0; i < totalDays; i++ {
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
