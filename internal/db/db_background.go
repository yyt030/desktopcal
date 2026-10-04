package db

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

func (d *DB) GetDayBackgroundsByDateRange(start, end string) (map[string]string, error) {
	return d.scanBackgrounds(
		`SELECT date, color FROM day_backgrounds WHERE date >= ? AND date <= ?`, start, end)
}

func (d *DB) GetAllBackgrounds() (map[string]string, error) {
	return d.scanBackgrounds(`SELECT date, color FROM day_backgrounds`)
}

// scanBackgrounds 把日期背景色查询结果统一收敛为 map[日期]颜色。
func (d *DB) scanBackgrounds(query string, args ...any) (map[string]string, error) {
	rows, err := d.conn.Query(query, args...)
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
