package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"desktopcal/internal/model"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New() (*DB, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	dbDir := filepath.Join(dir, "desktopcal")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	conn, err := sql.Open("sqlite3", filepath.Join(dbDir, "calendar.db"))
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}
	return db, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) migrate() error {
	_, err := d.conn.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			date TEXT NOT NULL,
			start_time TEXT DEFAULT '',
			end_time TEXT DEFAULT '',
			color TEXT DEFAULT '#4A90D9',
			completed INTEGER DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);

		CREATE TABLE IF NOT EXISTS day_backgrounds (
			date TEXT PRIMARY KEY,
			color TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	d.conn.Exec(`ALTER TABLE events ADD COLUMN completed INTEGER DEFAULT 0`)
	return nil
}

func (d *DB) GetEventsByDate(date string) ([]*model.Event, error) {
	rows, err := d.conn.Query(
		`SELECT id, title, description, date, start_time, end_time, color, completed, created_at, updated_at
		 FROM events WHERE date = ? ORDER BY completed ASC, start_time`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		e := &model.Event{}
		var completed int
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date,
			&e.StartTime, &e.EndTime, &e.Color, &completed, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.Completed = completed == 1
		events = append(events, e)
	}
	if events == nil {
		events = []*model.Event{}
	}
	return events, rows.Err()
}

func (d *DB) GetEventsByMonth(year, month int) ([]*model.Event, error) {
	start := fmt.Sprintf("%04d-%02d-01", year, month)
	next := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, time.Local)
	end := next.AddDate(0, 0, -1).Format("2006-01-02")
	rows, err := d.conn.Query(
		`SELECT id, title, description, date, start_time, end_time, color, completed, created_at, updated_at
		 FROM events WHERE date >= ? AND date <= ? ORDER BY date, completed ASC, start_time`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		e := &model.Event{}
		var completed int
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date,
			&e.StartTime, &e.EndTime, &e.Color, &completed, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.Completed = completed == 1
		events = append(events, e)
	}
	if events == nil {
		events = []*model.Event{}
	}
	return events, rows.Err()
}

func (d *DB) CreateEvent(event *model.Event) (*model.Event, error) {
	now := time.Now().Format(time.RFC3339)
	event.ID = uuid.New().String()
	event.CreatedAt = now
	event.UpdatedAt = now

	completed := 0
	if event.Completed {
		completed = 1
	}

	_, err := d.conn.Exec(
		`INSERT INTO events (id, title, description, date, start_time, end_time, color, completed, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.Title, event.Description, event.Date,
		event.StartTime, event.EndTime, event.Color, completed, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (d *DB) UpdateEvent(event *model.Event) error {
	event.UpdatedAt = time.Now().Format(time.RFC3339)
	completed := 0
	if event.Completed {
		completed = 1
	}
	_, err := d.conn.Exec(
		`UPDATE events SET title=?, description=?, date=?, start_time=?, end_time=?, color=?, completed=?, updated_at=?
		 WHERE id=?`,
		event.Title, event.Description, event.Date,
		event.StartTime, event.EndTime, event.Color, completed, event.UpdatedAt, event.ID)
	return err
}

func (d *DB) DeleteEvent(id string) error {
	_, err := d.conn.Exec(`DELETE FROM events WHERE id=?`, id)
	return err
}

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

func (d *DB) GetAllEvents() ([]*model.Event, error) {
	rows, err := d.conn.Query(
		`SELECT id, title, description, date, start_time, end_time, color, completed, created_at, updated_at
		 FROM events ORDER BY date, start_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		e := &model.Event{}
		var completed int
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date,
			&e.StartTime, &e.EndTime, &e.Color, &completed, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.Completed = completed == 1
		events = append(events, e)
	}
	if events == nil {
		events = []*model.Event{}
	}
	return events, rows.Err()
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

func (d *DB) GetSetting(key string) (string, error) {
	var value string
	err := d.conn.QueryRow(`SELECT value FROM settings WHERE key=?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (d *DB) SetSetting(key, value string) error {
	_, err := d.conn.Exec(
		`INSERT INTO settings (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value=excluded.value`, key, value)
	return err
}
