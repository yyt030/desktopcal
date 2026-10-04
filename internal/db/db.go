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

// New 打开（必要时创建）应用数据库并完成建表迁移。
func New() (*DB, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.TempDir()
	}
	dbDir := filepath.Join(dir, "desktopcal")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据目录失败: %w", err)
	}

	// busy_timeout：并发写入时等待而非立刻返回 SQLITE_BUSY；
	// WAL：读写可并发，避免界面查询被写事务阻塞。
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000&_journal_mode=WAL",
		filepath.Join(dbDir, "calendar.db"))
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	// 桌面端单用户场景：串行化连接，彻底规避 SQLite 写锁竞争。
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

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

	// 兼容旧库：completed 列在早期版本中缺失，确认后才补建，
	// 避免每次启动都执行一次必然失败的 ALTER TABLE。
	exists, err := d.columnExists("events", "completed")
	if err != nil {
		return fmt.Errorf("检查表结构失败: %w", err)
	}
	if !exists {
		if _, err := d.conn.Exec(`ALTER TABLE events ADD COLUMN completed INTEGER DEFAULT 0`); err != nil {
			return fmt.Errorf("补充 completed 列失败: %w", err)
		}
	}
	return nil
}

// columnExists 通过 PRAGMA table_info 判断表中是否已存在指定列。
func (d *DB) columnExists(table, column string) (bool, error) {
	rows, err := d.conn.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid          int
			name         string
			columnType   string
			notNull      int
			defaultValue sql.NullString
			primaryKey   int
		)
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}
