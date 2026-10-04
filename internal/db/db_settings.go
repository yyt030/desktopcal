package db

import "strings"

// GetSettings 批量读取指定键的设置值，未设置的键不会出现在结果中。
func (d *DB) GetSettings(keys []string) (map[string]string, error) {
	result := make(map[string]string)
	if len(keys) == 0 {
		return result, nil
	}

	placeholders := "?" + strings.Repeat(",?", len(keys)-1)
	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}

	rows, err := d.conn.Query(
		`SELECT key, value FROM settings WHERE key IN (`+placeholders+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, rows.Err()
}

// SetSettings 在单个事务中批量写入设置，保证读写一致。
func (d *DB) SetSettings(settings map[string]string) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value=excluded.value`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, value := range settings {
		if _, err := stmt.Exec(key, value); err != nil {
			return err
		}
	}
	return tx.Commit()
}
