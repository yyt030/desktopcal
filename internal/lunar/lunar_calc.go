package lunar

import (
	"fmt"
	"time"
)

func (c *Calendar) daysFrom1900(year, month, day int) int {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	base := time.Date(1900, 1, 31, 0, 0, 0, 0, time.UTC)
	return int(t.Sub(base).Hours() / 24)
}

func (c *Calendar) lunarYearDays(year int) int {
	idx := year - 1900
	if idx < 0 || idx >= len(lunarInfo) {
		return 365
	}
	info := lunarInfo[idx]
	sum := 29 * 12
	for i := 0; i < 12; i++ {
		if info&(0x1000>>i) != 0 {
			sum++
		}
	}
	if c.leapMonth(year) > 0 {
		sum += c.leapMonthDays(year)
	}
	return sum
}

func (c *Calendar) lunarMonthDays(year, month int) int {
	idx := year - 1900
	if idx < 0 || idx >= len(lunarInfo) {
		return 30
	}
	if lunarInfo[idx]&(0x1000>>(month-1)) != 0 {
		return 30
	}
	return 29
}

func (c *Calendar) leapMonth(year int) int {
	idx := year - 1900
	if idx < 0 || idx >= len(lunarInfo) {
		return 0
	}
	return lunarInfo[idx] & 0xf
}

func (c *Calendar) leapMonthDays(year int) int {
	idx := year - 1900
	if idx < 0 || idx >= len(lunarInfo) {
		return 29
	}
	if lunarInfo[idx]&0x10000 != 0 {
		return 30
	}
	return 29
}

func (c *Calendar) ganZhiYear(year int) string {
	idx := (year - 4) % 60
	gan := ganNames[idx%10]
	zhi := zhiNames[idx%12]
	return fmt.Sprintf("%s%s", gan, zhi)
}

func (c *Calendar) animalYear(year int) string {
	idx := (year - 4) % 12
	return animalNames[idx]
}

func (c *Calendar) monthName(month int, isLeap bool) string {
	if month < 1 || month > 12 {
		return ""
	}
	prefix := ""
	if isLeap {
		prefix = "闰"
	}
	return prefix + monthNames[month-1] + "月"
}

func (c *Calendar) dayName(day int) string {
	if day < 1 || day > 30 {
		return ""
	}
	return dayNames[day-1]
}
