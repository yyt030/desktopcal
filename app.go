package main

import (
	"context"
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
