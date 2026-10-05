package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"desktopcal/internal/model"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// exportData 是导入导出使用的 JSON 结构，字段与导出文件格式一一对应。
type exportData struct {
	Events      []*model.Event    `json:"events"`
	Backgrounds map[string]string `json:"backgrounds"`
	ExportedAt  string            `json:"exportedAt"`
}

func (a *App) ExportData(year, month int) (string, error) {
	var events []*model.Event
	var err error
	
	if year > 0 && month > 0 {
		// 导出指定月份的数据
		startDate := fmt.Sprintf("%d-%02d-01", year, month)
		endDate := fmt.Sprintf("%d-%02d-31", year, month)
		events, err = a.db.GetEventsByDateRange(startDate, endDate)
		if err != nil {
			return "", fmt.Errorf("导出事件失败: %w", err)
		}
	} else {
		// 导出全部数据
		events, err = a.db.GetAllEvents()
		if err != nil {
			return "", fmt.Errorf("导出事件失败: %w", err)
		}
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
		return "", fmt.Errorf("序列化导出数据失败: %w", err)
	}
	return string(b), nil
}

func (a *App) ExportDataToFile(year, month int) error {
	monthStr := fmt.Sprintf("%d-%02d", year, month)
	defaultFilename := fmt.Sprintf("desktopcal-backup-%s.json", monthStr)

	settings := a.GetSettings()
	exportDir := settings.ExportPath
	if exportDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("未设置导出路径，且无法获取用户主目录")
		}
		exportDir = homeDir
	}

	savePath := filepath.Join(exportDir, defaultFilename)

	jsonData, err := a.ExportData(year, month)
	if err != nil {
		return err
	}

	if err := os.WriteFile(savePath, []byte(jsonData), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
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

func (a *App) PickDirectory() (string, error) {
	if a.app == nil || a.app.Dialog == nil {
		return "", fmt.Errorf("对话框不可用")
	}

	dialog := a.app.Dialog.OpenFile()
	dialog.SetMessage("选择导出目录")
	dialog.SetOptions(&application.OpenFileDialogOptions{
		CanChooseDirectories: true,
		CanChooseFiles:       false,
	})

	dir, err := dialog.PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("打开目录选择对话框失败: %w", err)
	}

	return dir, nil
}
