package lunar

import (
	"time"

	"desktopcal/internal/model"
)

type Calendar struct{}

func NewCalendar() *Calendar {
	return &Calendar{}
}

func (c *Calendar) SolarToLunar(t time.Time) *model.LunarDay {
	year := t.Year()
	month := int(t.Month())
	day := t.Day()

	offset := c.daysFrom1900(year, month, day)

	lunarYear := 1900
	for i := 1900; i < 2101; i++ {
		yd := c.lunarYearDays(i)
		if offset < yd {
			break
		}
		offset -= yd
		lunarYear = i + 1
	}
	if lunarYear > 2100 {
		lunarYear = 2100
	}

	leap := c.leapMonth(lunarYear)
	isLeap := false
	lunarMonth := 1

	for i := 1; i <= 12; i++ {
		md := c.lunarMonthDays(lunarYear, i)
		if offset < md {
			lunarMonth = i
			break
		}
		offset -= md

		if i == leap {
			lmd := c.leapMonthDays(lunarYear)
			if offset < lmd {
				isLeap = true
				lunarMonth = i
				break
			}
			offset -= lmd
		}
	}

	lunarDay := offset + 1
	if lunarDay < 1 {
		lunarDay = 1
	}
	if lunarDay > 30 {
		lunarDay = 30
	}

	return &model.LunarDay{
		Year:        lunarYear,
		Month:       lunarMonth,
		Day:         lunarDay,
		IsLeapMonth: isLeap,
		DayName:     c.dayName(lunarDay),
		MonthName:   c.monthName(lunarMonth, isLeap),
		Animal:      c.animalYear(lunarYear),
		GanZhi:      c.ganZhiYear(lunarYear),
	}
}

func (c *Calendar) GetSolarTerm(t time.Time) string {
	month := int(t.Month())
	day := t.Day()

	for i := 0; i < 24; i++ {
		if solarTermDays[i][0] == month && solarTermDays[i][1] == day {
			return solarTermNames[i]
		}
	}
	return ""
}
