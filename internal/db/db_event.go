package db

import (
	"database/sql"
	"time"

	"desktopcal/internal/model"

	"github.com/google/uuid"
)

const eventColumns = `id, title, description, date, start_time, end_time, color, completed, created_at, updated_at`

func scanEvent(rows *sql.Rows) (*model.Event, error) {
	e := &model.Event{}
	var completed int
	if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date,
		&e.StartTime, &e.EndTime, &e.Color, &completed, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return nil, err
	}
	e.Completed = completed == 1
	return e, nil
}

func scanEvents(rows *sql.Rows) ([]*model.Event, error) {
	defer rows.Close()
	events := []*model.Event{}
	for rows.Next() {
		e, err := scanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (d *DB) GetEventsByDate(date string) ([]*model.Event, error) {
	rows, err := d.conn.Query(
		`SELECT `+eventColumns+` FROM events WHERE date = ?
		 ORDER BY completed ASC, start_time`, date)
	if err != nil {
		return nil, err
	}
	return scanEvents(rows)
}

func (d *DB) GetEventsByDateRange(start, end string) ([]*model.Event, error) {
	rows, err := d.conn.Query(
		`SELECT `+eventColumns+` FROM events WHERE date >= ? AND date <= ?
		 ORDER BY date, completed ASC, start_time`, start, end)
	if err != nil {
		return nil, err
	}
	return scanEvents(rows)
}

func (d *DB) GetAllEvents() ([]*model.Event, error) {
	rows, err := d.conn.Query(
		`SELECT ` + eventColumns + ` FROM events ORDER BY date, start_time`)
	if err != nil {
		return nil, err
	}
	return scanEvents(rows)
}

func (d *DB) CreateEvent(event *model.Event) (*model.Event, error) {
	now := time.Now().Format(time.RFC3339)
	event.ID = uuid.New().String()
	event.CreatedAt = now
	event.UpdatedAt = now

	_, err := d.conn.Exec(
		`INSERT INTO events (`+eventColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.Title, event.Description, event.Date,
		event.StartTime, event.EndTime, event.Color, boolToInt(event.Completed),
		event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (d *DB) UpdateEvent(event *model.Event) error {
	event.UpdatedAt = time.Now().Format(time.RFC3339)
	_, err := d.conn.Exec(
		`UPDATE events SET title=?, description=?, date=?, start_time=?, end_time=?,
			color=?, completed=?, updated_at=? WHERE id=?`,
		event.Title, event.Description, event.Date,
		event.StartTime, event.EndTime, event.Color, boolToInt(event.Completed),
		event.UpdatedAt, event.ID)
	return err
}

func (d *DB) DeleteEvent(id string) error {
	_, err := d.conn.Exec(`DELETE FROM events WHERE id=?`, id)
	return err
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
