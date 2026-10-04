package main

import (
	"fmt"
	"time"

	"desktopcal/internal/model"
)

// 单日事件的增删改查、日期背景色与今日农历/节气查询。

func (a *App) GetEventsByDate(date string) ([]*model.Event, error) {
	return a.db.GetEventsByDate(date)
}

func (a *App) CreateEvent(event *model.Event) (*model.Event, error) {
	if event == nil {
		return nil, fmt.Errorf("事件不能为空")
	}
	if event.Title == "" {
		return nil, fmt.Errorf("事件标题不能为空")
	}
	if event.Date == "" {
		return nil, fmt.Errorf("事件日期不能为空")
	}
	return a.db.CreateEvent(event)
}

func (a *App) UpdateEvent(event *model.Event) error {
	if event == nil || event.ID == "" {
		return fmt.Errorf("事件 ID 不能为空")
	}
	if event.Title == "" {
		return fmt.Errorf("事件标题不能为空")
	}
	return a.db.UpdateEvent(event)
}

func (a *App) DeleteEvent(id string) error {
	if id == "" {
		return fmt.Errorf("事件 ID 不能为空")
	}
	return a.db.DeleteEvent(id)
}

func (a *App) SetDayBackground(date, color string) error {
	if date == "" {
		return fmt.Errorf("日期不能为空")
	}
	return a.db.SetDayBackground(date, color)
}

func (a *App) GetTodayLunar() *model.LunarDay {
	return a.lcal.SolarToLunar(time.Now())
}

func (a *App) GetTodaySolarTerm() string {
	return a.lcal.GetSolarTerm(time.Now())
}
