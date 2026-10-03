package db

import (
	"desktopcal/internal/model"
)

func (d *DB) ImportEvents(events []*model.Event) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM events`)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO events (id, title, description, date, start_time, end_time, color, completed, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range events {
		completed := 0
		if e.Completed {
			completed = 1
		}
		_, err := stmt.Exec(e.ID, e.Title, e.Description, e.Date,
			e.StartTime, e.EndTime, e.Color, completed, e.CreatedAt, e.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (d *DB) ImportBackgrounds(backgrounds map[string]string) error {
	tx, err := d.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM day_backgrounds`)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO day_backgrounds (date, color) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for date, color := range backgrounds {
		if _, err := stmt.Exec(date, color); err != nil {
			return err
		}
	}
	return tx.Commit()
}
