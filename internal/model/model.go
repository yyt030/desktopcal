package model

// DateLayout 是日期在数据库、后端与前端之间统一使用的格式。
const DateLayout = "2006-01-02"

// 日历显示模式。
const (
	LayoutCompact = "compact" // 月历连续排布，交接处无空格
	LayoutLoose   = "loose"   // 每月单独一页，不混入其它月
)

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Color       string `json:"color"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type LunarDay struct {
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	IsLeapMonth bool   `json:"isLeapMonth"`
	DayName     string `json:"dayName"`
	MonthName   string `json:"monthName"`
	Animal      string `json:"animal"`
	GanZhi      string `json:"ganZhi"`
}

type DayInfo struct {
	SolarDate string    `json:"solarDate"`
	Lunar     *LunarDay `json:"lunar"`
	SolarTerm string    `json:"solarTerm"`
	Events    []*Event  `json:"events"`
	IsToday   bool      `json:"isToday"`
	IsCurrent bool      `json:"isCurrentMonth"`
	IsEmpty   bool      `json:"isEmpty"`
	BgColor   string    `json:"bgColor"`
}

type Settings struct {
	Opacity            int    `json:"opacity"`
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	Theme              string `json:"theme"`
	AlwaysTop          bool   `json:"alwaysTop"`
	FontSize           int    `json:"fontSize"`
	FontFamily         string `json:"fontFamily"`
	CalendarFontSize   int    `json:"calendarFontSize"`
	CalendarFontFamily string `json:"calendarFontFamily"`
	FirstDayOfWeek     int    `json:"firstDayOfWeek"`
	ShowLunar          bool   `json:"showLunar"`
	ShowSolarTerm      bool   `json:"showSolarTerm"`
	LayoutMode         string `json:"layoutMode"`
	Hotkey             string `json:"hotkey"`
	CloseToTray        bool   `json:"closeToTray"`
	AutoStart          bool   `json:"autoStart"`
	ExportPath         string `json:"exportPath"`
	EditorWidth        int    `json:"editorWidth"`
	EditorHeight       int    `json:"editorHeight"`
}

// DefaultSettings 返回设置项的默认值：既用于后端读取失败时兜底，
// 也是前端默认值与设置项编解码（见 app_settings.go）的单一来源。
func DefaultSettings() *Settings {
	return &Settings{
		Opacity:            100,
		Width:              1100,
		Height:             750,
		Theme:              "dark",
		AlwaysTop:          true,
		FontSize:           13,
		FontFamily:         "",
		CalendarFontSize:   14,
		CalendarFontFamily: "",
		FirstDayOfWeek:     1,
		ShowLunar:          true,
		ShowSolarTerm:      true,
		LayoutMode:         LayoutLoose,
		Hotkey:             DefaultHotkey(),
		CloseToTray:        true,
		ExportPath:         "",
		EditorWidth:        420,
		EditorHeight:       400,
	}
}

// DefaultHotkey 返回平台相关的默认快捷键。
func DefaultHotkey() string {
	// macOS 用 Cmd+Shift+C，其它平台用 Ctrl+Shift+C
	if isMacOS() {
		return "cmd+shift+c"
	}
	return "ctrl+shift+c"
}

func isMacOS() bool {
	// 编译时判断，避免引入 runtime 包
	return darwin
}
