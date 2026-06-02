package logger

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/d3v-sh/sec_agent/config"
	_ "modernc.org/sqlite"
)

type ToolLog struct {
	ID        int       `db:"id"`
	SessionID string    `db:"session_id"`
	Time      time.Time `db:"time"`
	Tool      string    `db:"tool"`
	Args      string    `db:"args"`
	Result    string    `db:"result"`
}
type Session struct {
	ID        string    `db:"id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

var (
	db      *sql.DB
	current *config.Session
)

func Current() *config.Session {
	return current
}

func Start() error {
	os.MkdirAll("data", 0755)

	var err error
	db, err = sql.Open("sqlite", "data/agent.db")
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions(id TEXT PRIMARY KEY, start_time DATETIME, end_time DATETIME);
		CREATE TABLE IF NOT EXISTS tool_logs(id INTEGER PRIMARY KEY AUTOINCREMENT, session_id TEXT, time DATETIME, tool TEXT, args TEXT, result TEXT, FOREIGN KEY(session_id) REFERENCES sessions(id));
		CREATE TABLE IF NOT EXISTS targets (host TEXT PRIMARY KEY, info TEXT, updated DATETIME);
		`)
	if err != nil {
		return err
	}

	current = &config.Session{
		ID:        time.Now().Format("20060102-150405"),
		StartTime: time.Now(),
	}
	_, err = db.Exec(
		"INSERT INTO sessions (id, start_time) VALUES (?, ?)",
		current.ID, current.StartTime,
	)
	return err
}

func LogTool(tool, args, result string) {
	if db == nil || current == nil {
		return
	}
	db.Exec(
		"INSERT INTO tool_logs (session_id, time, tool, args, result) VALUES (?, ?, ?, ?, ?)",
		current.ID, time.Now(), tool, args, result,
	)
}
func End() {
	if db == nil || current == nil {
		return
	}
	db.Exec(
		"UPDATE sessions SET end_time = ? WHERE id = ?",
		time.Now(), current.ID,
	)
	fmt.Printf("\nSession %s saved to data/agent.db\n", current.ID)
	db.Close()
}
func DB() *sql.DB {
	return db
}
