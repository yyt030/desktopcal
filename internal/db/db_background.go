package db

import (
	"fmt"
	"time"
)

func (d *DB) SetDayBackground(date, color string) error {
	if color == "" {
		_, err := d.conn.Exec(`DELETE FROM day_backgrounds WHERE date=?`, date)
		return err
	}
	_, err := d.conn.Exec(
		`INSERT INTO day_backgrounds (date, color) VALUES (?, ?)
		 ON CONFLICT(date) DO UPDATE SET color=excluded.color`, date, color)
	return err
}

func (d *DB) GetDayBackgroundsByMonth(year, month int) (map[string]string, error) {
	start := fmt.Sprintf("%04d-%02d-01", year, month)
	next := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.Local)
	end := next.AddDate(0, 0, -1).Format("2006-01-02")
	rows, err := d.conn.Query(
		`SELECT date, color FROM day_backgrounds WHERE date >= ? AND date <= ?`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var date, color string
		if err := rows.Scan(&date, &color); err != nil {
			return nil, err
		}
		result[date] = color
	}
	return result, rows.Err()
}

func (d *DB) GetDayBackgroundsByDateRange(start, end string) (map[string]string, error) {
	rows, err := d.conn.Query(
		`SELECT date, color FROM day_backgrounds WHERE date >= ? AND date <= ?`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var date, color string
		if err := rows.Scan(&date, &color); err != nil {
			return nil, err
		}
		result[date] = color
	}
	return result, rows.Err()
}

func (d *DB) GetAllBackgrounds() (map[string]string, error) {
	rows, err := d.conn.Query(`SELECT date, color FROM day_backgrounds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var date, color string
		if err := rows.Scan(&date, &color); err != nil {
			return nil, err
		}
		result[date] = color
	}
	return result, rows.Err()
}
