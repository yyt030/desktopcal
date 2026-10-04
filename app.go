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

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, -1)

	layoutMode := "compact"
	if v, _ := a.db.GetSetting("layout_mode"); v != "" {
		layoutMode = v
	}
	firstDayOfWeek := 0
	if v, _ := a.db.GetSetting("first_day_of_week"); v != "" {
		fmt.Sscanf(v, "%d", &firstDayOfWeek)
	}

	// 按星期列对齐：1号前的空格数
	offset := (int(start.Weekday()) - firstDayOfWeek + 7) % 7
	monthDays := end.Day()

	gridStart := start.AddDate(0, 0, -offset)
	emptyBefore := 0
	totalDays := 0
	isCurrent := today.Year() == year && today.Month() == time.Month(month)

	switch {
	case layoutMode == "loose":
		// 宽松模式：每月单独一页，1号前和月末不足部分留空格
		totalDays = ((offset + monthDays + 6) / 7) * 7
		emptyBefore = offset
	case isCurrent:
		// 紧凑模式当月：今天所在的周居中（上方2行、下方1行），共4周
		todayOffset := (int(today.Weekday()) - firstDayOfWeek + 7) % 7
		gridStart = today.AddDate(0, 0, -todayOffset-14)
		totalDays = 28
	default:
		// 紧凑模式其它月：从1号所在的周到月末所在的周，相邻月日期连续补齐，无空格
		totalDays = ((offset + monthDays + 6) / 7) * 7
	}
	emptyAfter := totalDays - emptyBefore - monthDays
	if layoutMode != "loose" {
		emptyAfter = 0
	}

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
		if i < emptyBefore || i >= totalDays-emptyAfter {
			// 宽松模式首尾补位空格
			days[i] = &model.DayInfo{
				SolarDate: "",
				Events:    []*model.Event{},
				IsEmpty:   true,
			}
			continue
		}
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
