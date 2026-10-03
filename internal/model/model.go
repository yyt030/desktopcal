package model

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

type SolarTerm struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type DayInfo struct {
	SolarDate string    `json:"solarDate"`
	Lunar     *LunarDay `json:"lunar"`
	SolarTerm string    `json:"solarTerm"`
	Events    []*Event  `json:"events"`
	IsToday   bool      `json:"isToday"`
	IsCurrent bool      `json:"isCurrentMonth"`
	BgColor   string    `json:"bgColor"`
}

type Settings struct {
	Opacity       int    `json:"opacity"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Theme         string `json:"theme"`
	AlwaysTop     bool   `json:"alwaysTop"`
	FontSize      int    `json:"fontSize"`
	FirstDayOfWeek int   `json:"firstDayOfWeek"`
}
