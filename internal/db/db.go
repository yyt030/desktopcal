package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
