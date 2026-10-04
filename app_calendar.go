package main

import (
	"strconv"
	"time"

	"desktopcal/internal/model"
)

// GetMonthDays 返回渲染月视图所需的日期格子，包含农历、节气、事件与背景色。
// 格子数量随显示模式变化：宽松模式补齐整周并留空首尾，紧凑模式连续排布。
func (a *App) GetMonthDays(year, month int) ([]*model.DayInfo, error) {
	today := time.Now()
	todayStr := today.Format(model.DateLayout)

	layoutMode, firstDayOfWeek := a.calendarSettings()

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	monthDays := start.AddDate(0, 1, -1).Day()

	// 按星期列对齐：1 号之前需要留出的占位格数
	offset := (int(start.Weekday()) - firstDayOfWeek + 7) % 7
	gridStart := start.AddDate(0, 0, -offset)
	emptyBefore := 0
	isCurrent := today.Year() == year && today.Month() == time.Month(month)

	var totalDays int
	switch {
	case layoutMode == model.LayoutLoose:
		// 宽松模式：每月单独一页，首尾补空对齐整周
		totalDays = weekAlignedGrid(offset + monthDays)
		emptyBefore = offset
	case isCurrent:
		// 紧凑模式当月：今天所在的周居中（上方 2 行、下方 1 行），共 4 周
		todayOffset := (int(today.Weekday()) - firstDayOfWeek + 7) % 7
		gridStart = today.AddDate(0, 0, -todayOffset-14)
		totalDays = 28
	default:
		// 紧凑模式其它月：从 1 号所在周排到月末所在周，相邻月日期连续补齐，无空格
		totalDays = weekAlignedGrid(offset + monthDays)
	}

	emptyAfter := 0
	if layoutMode == model.LayoutLoose {
		emptyAfter = totalDays - emptyBefore - monthDays
	}

	gridStartStr := gridStart.Format(model.DateLayout)
	gridEndStr := gridStart.AddDate(0, 0, totalDays-1).Format(model.DateLayout)
	eventsByDate := a.eventsByDate(gridStartStr, gridEndStr)
	bgColors := a.dayBackgrounds(gridStartStr, gridEndStr)

	days := make([]*model.DayInfo, totalDays)
	for i := 0; i < totalDays; i++ {
		if i < emptyBefore || i >= totalDays-emptyAfter {
			days[i] = &model.DayInfo{Events: []*model.Event{}, IsEmpty: true}
			continue
		}
		d := gridStart.AddDate(0, 0, i)
		dateStr := d.Format(model.DateLayout)
		events := eventsByDate[dateStr]
		if events == nil {
			events = []*model.Event{}
		}
		days[i] = &model.DayInfo{
			SolarDate: dateStr,
			Lunar:     a.lcal.SolarToLunar(d),
			SolarTerm: a.lcal.GetSolarTerm(d),
			Events:    events,
			IsToday:   dateStr == todayStr,
			IsCurrent: d.Year() == year && d.Month() == time.Month(month),
			BgColor:   bgColors[dateStr],
		}
	}
	return days, nil
}

// weekAlignedGrid 把天数向上补齐到整周（7 的倍数）。
func weekAlignedGrid(days int) int {
	return ((days + 6) / 7) * 7
}

// calendarSettings 读取月视图渲染所需的最小设置，缺省或非法值回退到默认值。
func (a *App) calendarSettings() (layoutMode string, firstDayOfWeek int) {
	defaults := model.DefaultSettings()
	layoutMode = defaults.LayoutMode
	firstDayOfWeek = defaults.FirstDayOfWeek

	stored, err := a.db.GetSettings(calendarSettingKeys)
	if err != nil {
		return layoutMode, firstDayOfWeek
	}
	if v := stored["layout_mode"]; v != "" {
		layoutMode = v
	}
	if v, err := strconv.Atoi(stored["first_day_of_week"]); err == nil {
		firstDayOfWeek = v
	}
	return layoutMode, firstDayOfWeek
}

// eventsByDate 查询区间内事件并按日期分组；查询失败时返回空集合，不影响日历渲染。
func (a *App) eventsByDate(start, end string) map[string][]*model.Event {
	byDate := make(map[string][]*model.Event)
	events, err := a.db.GetEventsByDateRange(start, end)
	if err != nil {
		return byDate
	}
	for _, e := range events {
		byDate[e.Date] = append(byDate[e.Date], e)
	}
	return byDate
}

// dayBackgrounds 查询区间内的日期背景色；查询失败时返回空集合。
func (a *App) dayBackgrounds(start, end string) map[string]string {
	colors, err := a.db.GetDayBackgroundsByDateRange(start, end)
	if err != nil {
		return make(map[string]string)
	}
	return colors
}
