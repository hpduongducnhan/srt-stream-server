package app

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db              *sql.DB
	dbMutex         sync.RWMutex
	allowedIPs      []string
	allowedStreams  []AppStreamPair
	filterDataMutex sync.RWMutex
)

// AppStreamPair represents an allowed app/stream combination
type AppStreamPair struct {
	App    string
	Stream string
}

// InitDatabase initializes the SQLite database and creates tables if not exist
func InitDatabase(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Create tables if not exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS allowed_ips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL UNIQUE,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS allowed_streams (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app TEXT NOT NULL,
		stream TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(app, stream)
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	// Load initial data
	return ReloadFilterData()
}

// ReloadFilterData loads allowed IPs and streams from database into memory
func ReloadFilterData() error {
	filterDataMutex.Lock()
	defer filterDataMutex.Unlock()

	// Load allowed IPs
	ips, err := loadAllowedIPsFromDB()
	if err != nil {
		return err
	}
	allowedIPs = ips

	// Load allowed streams
	streams, err := loadAllowedStreamsFromDB()
	if err != nil {
		return err
	}
	allowedStreams = streams

	logger.Info().
		Int("allowed_ips_count", len(allowedIPs)).
		Int("allowed_streams_count", len(allowedStreams)).
		Msg("filter data reloaded from database")

	return nil
}

func loadAllowedIPsFromDB() ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	rows, err := db.Query("SELECT ip FROM allowed_ips")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, err
		}
		ips = append(ips, ip)
	}
	return ips, rows.Err()
}

func loadAllowedStreamsFromDB() ([]AppStreamPair, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	rows, err := db.Query("SELECT app, stream FROM allowed_streams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var streams []AppStreamPair
	for rows.Next() {
		var pair AppStreamPair
		if err := rows.Scan(&pair.App, &pair.Stream); err != nil {
			return nil, err
		}
		streams = append(streams, pair)
	}
	return streams, rows.Err()
}

// GetAllowedIPs returns a copy of allowed IPs
func GetAllowedIPs() []string {
	filterDataMutex.RLock()
	defer filterDataMutex.RUnlock()

	result := make([]string, len(allowedIPs))
	copy(result, allowedIPs)
	return result
}

// GetAllowedStreams returns a copy of allowed streams
func GetAllowedStreams() []AppStreamPair {
	filterDataMutex.RLock()
	defer filterDataMutex.RUnlock()

	result := make([]AppStreamPair, len(allowedStreams))
	copy(result, allowedStreams)
	return result
}

// AddAllowedIP adds an IP to the allowed list
func AddAllowedIP(ip string, description string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec("INSERT OR IGNORE INTO allowed_ips (ip, description) VALUES (?, ?)", ip, description)
	if err != nil {
		return err
	}

	return ReloadFilterData()
}

// RemoveAllowedIP removes an IP from the allowed list
func RemoveAllowedIP(ip string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec("DELETE FROM allowed_ips WHERE ip = ?", ip)
	if err != nil {
		return err
	}

	return ReloadFilterData()
}

// AddAllowedStream adds an app/stream pair to the allowed list
func AddAllowedStream(app, stream, description string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec("INSERT OR IGNORE INTO allowed_streams (app, stream, description) VALUES (?, ?, ?)", app, stream, description)
	if err != nil {
		return err
	}

	return ReloadFilterData()
}

// RemoveAllowedStream removes an app/stream pair from the allowed list
func RemoveAllowedStream(app, stream string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	_, err := db.Exec("DELETE FROM allowed_streams WHERE app = ? AND stream = ?", app, stream)
	if err != nil {
		return err
	}

	return ReloadFilterData()
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
