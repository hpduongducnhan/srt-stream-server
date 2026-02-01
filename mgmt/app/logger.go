package app

import (
	"os"

	"github.com/rs/zerolog"
)

func initLogger() zerolog.Logger {
	// logPath := os.Getenv("NDD_LOG_PATH")
	// if logPath == "" {
	// 	logPath, _ = os.Getwd()
	// }
	// if logPath != "" {
	// 	// Ensure the directory exists
	// 	if err := os.MkdirAll(logPath, os.ModePerm); err != nil {
	// 		panic("failed to create log directory: " + err.Error())
	// 	}
	// }
	// writer := &lumberjack.Logger{
	// 	Filename:   logPath + "/nddsrt_app.log",
	// 	MaxSize:    100, // MB
	// 	MaxBackups: 7,
	// 	MaxAge:     14, // days
	// 	Compress:   true,
	// }

	// multi := zerolog.MultiLevelWriter(os.Stdout, writer)
	multi := zerolog.MultiLevelWriter(os.Stdout)
	return zerolog.New(multi).
		With().
		Timestamp().
		Logger()
}

func GetLogger() zerolog.Logger {
	return logger
}

func init() {
	logger = initLogger()
}
