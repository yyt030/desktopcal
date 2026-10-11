package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"desktopcal/internal/model"
)

var (
	textDateHeader = regexp.MustCompile(`^\[(\d{4}-\d{2}-\d{2})\]$`)
	textTaskPrefix = regexp.MustCompile(`^\d+\.\s*`)
)

// parseTextEvents 解析桌面日历官方软件导出的文本格式：
// [YYYY-MM-DD] 日期头 + 虚线分隔 + “1. 任务” 编号行，按日期分组返回。
// 行首的 ✓/√ 前缀视为已完成；无编号的续行并入上一条任务。
func parseTextEvents(text string) (map[string][]*model.Event, error) {
	events := map[string][]*model.Event{}
	currentDate := ""

	for _, raw := range strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.Trim(line, "-") == "" {
			continue
		}
		if m := textDateHeader.FindStringSubmatch(line); m != nil {
			if _, err := time.Parse("2006-01-02", m[1]); err != nil {
				return nil, fmt.Errorf("无效的日期头：%s", m[1])
			}
			currentDate = m[1]
			continue
		}
		if currentDate == "" {
			continue
		}

		dayEvents := events[currentDate]
		if !textTaskPrefix.MatchString(line) && len(dayEvents) > 0 {
			last := dayEvents[len(dayEvents)-1]
			last.Title += "\n" + line
			continue
		}
		body := strings.TrimSpace(textTaskPrefix.ReplaceAllString(line, ""))
		completed := false
		for _, mark := range []string{"✓", "√"} {
			if strings.HasPrefix(body, mark) {
				completed = true
				body = strings.TrimSpace(strings.TrimPrefix(body, mark))
			}
		}
		for _, suffix := range []string{"(已完成)", "（已完成）"} {
			if strings.HasSuffix(body, suffix) {
				completed = true
				body = strings.TrimSpace(strings.TrimSuffix(body, suffix))
			}
		}
		if body == "" {
			continue
		}
		events[currentDate] = append(dayEvents, &model.Event{
			Title:     body,
			Date:      currentDate,
			Color:     "#4A90D9",
			Completed: completed,
		})
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("文本中未找到任何日期数据")
	}
	return events, nil
}

// importTextData 合并导入：按 (date, title) 去重，已存在的任务跳过，仅插入新任务。
func (a *App) importTextData(text string) error {
	events, err := parseTextEvents(text)
	if err != nil {
		return err
	}
	for date, dayEvents := range events {
		existing, err := a.db.GetEventsByDate(date)
		if err != nil {
			return fmt.Errorf("查询 %s 的已有事件失败: %w", date, err)
		}
		existingTitles := make(map[string]bool, len(existing))
		for _, e := range existing {
			existingTitles[e.Title] = true
		}
		var toInsert []*model.Event
		for _, e := range dayEvents {
			if !existingTitles[e.Title] {
				toInsert = append(toInsert, e)
			}
		}
		if len(toInsert) == 0 {
			continue
		}
		if err := a.db.BatchSaveEvents(date, append(existing, toInsert...)); err != nil {
			return fmt.Errorf("导入 %s 的事件失败: %w", date, err)
		}
	}
	return nil
}

// ImportData 自动识别导入格式：JSON 备份整体覆盖，桌面日历官方文本按日期覆盖。
func (a *App) ImportData(content string) error {
	var data exportData
	if err := json.Unmarshal([]byte(content), &data); err == nil {
		if err := a.db.ImportEvents(data.Events); err != nil {
			return fmt.Errorf("导入事件失败: %w", err)
		}
		if err := a.db.ImportBackgrounds(data.Backgrounds); err != nil {
			return fmt.Errorf("导入背景失败: %w", err)
		}
		return nil
	}
	return a.importTextData(content)
}
