package app

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var enableStreamFilter bool
var filterMutex sync.RWMutex

func init() {
	logger = initLogger()
	envStreamFilter := os.Getenv("ENABLE_STREAM_FILTER")
	if envStreamFilter == "" {
		enableStreamFilter = true
	} else {
		if envStreamFilter == "1" || envStreamFilter == "true" || envStreamFilter == "TRUE" {
			enableStreamFilter = true
		} else {
			enableStreamFilter = false
		}
	}

	// Initialize database
	dbPath := os.Getenv("NDD_DB_PATH")
	if dbPath == "" {
		// Default database path
		execPath, err := os.Executable()
		if err != nil {
			logger.Error().Err(err).Msg("failed to get executable path")
			dbPath = "./stream-server.db"
		} else {
			dbPath = filepath.Join(filepath.Dir(execPath), "stream-server.db")
		}
	}

	if err := InitDatabase(dbPath); err != nil {
		logger.Error().Err(err).Str("db_path", dbPath).Msg("failed to initialize filter database")
	} else {
		logger.Info().Str("db_path", dbPath).Msg("filter database initialized")
	}
}

// SetStreamFilterEnabled enables or disables stream filtering
func SetStreamFilterEnabled(enabled bool) {
	filterMutex.Lock()
	defer filterMutex.Unlock()
	enableStreamFilter = enabled
	logger.Info().Bool("enabled", enabled).Msg("stream filter status changed")
}

// GetStreamFilterEnabled returns current stream filter status
func GetStreamFilterEnabled() bool {
	filterMutex.RLock()
	defer filterMutex.RUnlock()
	return enableStreamFilter
}
