package logger

import (
	"os"
	"time"
)

func InitLogger() {
	loggerConfig := LoggerConfig{
		Mode:        os.Getenv("MODE"),
		LogDir:      "log",
		MaxSizeMB:   100,
		MaxBackups:  10,
		MaxAgeDays:  30,
		Compress:    true,
		Console:     true,
		ShowCaller:  true,
		TimeZone:    "Asia/Shanghai",
		LogFileName: time.Now().Format("2006-01-02") + ".log",
	}

	Init(&loggerConfig)
}
