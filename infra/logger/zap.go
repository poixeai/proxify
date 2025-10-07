package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	Mode        string // "debug" / "release"
	LogDir      string // directory for log files
	MaxSizeMB   int    // max size per log file in MB
	MaxBackups  int    // max number of backup files
	MaxAgeDays  int    // keep log files for N days
	Compress    bool   // if true, compress rotated log files
	Console     bool   // if true, log to console
	ShowCaller  bool   // if true, show caller info
	TimeZone    string // time zone for timestamps
	LogFileName string // log file name
}

// global logger
var ZapLog *zap.SugaredLogger

// Init initializes the global logger
func Init(config *LoggerConfig) {
	if config == nil {
		config = defaultLoggerConfig()
	}
	logger := newLogger(config)
	ZapLog = logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
	defer ZapLog.Sync()
}

// default configuration
func defaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Mode:        "debug",
		LogDir:      "log",
		MaxSizeMB:   500,
		MaxBackups:  20,
		MaxAgeDays:  14,
		Compress:    false,
		Console:     true,
		ShowCaller:  true,
		TimeZone:    "Asia/Shanghai",
		LogFileName: time.Now().Format("2006-01-02") + ".log",
	}
}

// construct a new zap.Logger based on the config
func newLogger(cfg *LoggerConfig) *zap.Logger {
	level := zapcore.InfoLevel
	if cfg.Mode == "debug" {
		level = zapcore.DebugLevel
	}

	// create log directory if not exists
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		fmt.Printf("ERROR: unable to create log dir: %v\n", err)
		os.Exit(1)
	}

	// create multiple cores
	var cores []zapcore.Core

	if cfg.Console {
		consoleCore := zapcore.NewCore(
			getConsoleEncoder(cfg.TimeZone),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	filePath := filepath.Join(cfg.LogDir, cfg.LogFileName)
	fileCore := zapcore.NewCore(
		getFileEncoder(cfg.TimeZone),
		zapcore.AddSync(getLumberjackWriter(filePath, cfg)),
		level,
	)
	cores = append(cores, fileCore)

	core := zapcore.NewTee(cores...)

	options := []zap.Option{}
	if cfg.ShowCaller {
		options = append(options, zap.AddCaller())
	}

	return zap.New(core, options...)
}

// console output encoder (human-readable)
func getConsoleEncoder(tz string) zapcore.Encoder {
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encCfg.EncodeTime = makeTimeEncoder(tz)
	encCfg.ConsoleSeparator = " | "
	return zapcore.NewConsoleEncoder(encCfg)
}

// file output encoder (JSON)
func getFileEncoder(tz string) zapcore.Encoder {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encCfg.EncodeTime = makeTimeEncoder(tz)
	return zapcore.NewJSONEncoder(encCfg)
}

// time encoder with timezone
func makeTimeEncoder(tz string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		loc, _ := time.LoadLocation(tz)
		enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
	}
}

// file writer with rotation
func getLumberjackWriter(filename string, cfg *LoggerConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	})
}

/* --------------------- Wrapper ---------------------- */

// Level gets the current log level
func Level() zapcore.Level {
	return ZapLog.Level()
}

// common log methods
func Log(lvl zapcore.Level, args ...interface{}) { ZapLog.Log(lvl, args...) }
func Logw(lvl zapcore.Level, msg string, keysAndValues ...interface{}) {
	ZapLog.Logw(lvl, msg, keysAndValues...)
}
func Logf(lvl zapcore.Level, template string, args ...interface{}) {
	ZapLog.Logf(lvl, template, args...)
}

// Debug
func Debug(args ...interface{})                   { ZapLog.Debug(args...) }
func Debugw(msg string, kv ...interface{})        { ZapLog.Debugw(msg, kv...) }
func Debugf(template string, args ...interface{}) { ZapLog.Debugf(template, args...) }

// Info
func Info(args ...interface{})                   { ZapLog.Info(args...) }
func Infow(msg string, kv ...interface{})        { ZapLog.Infow(msg, kv...) }
func Infof(template string, args ...interface{}) { ZapLog.Infof(template, args...) }

// Warn
func Warn(args ...interface{})                   { ZapLog.Warn(args...) }
func Warnw(msg string, kv ...interface{})        { ZapLog.Warnw(msg, kv...) }
func Warnf(template string, args ...interface{}) { ZapLog.Warnf(template, args...) }

// Error
func Error(args ...interface{})                   { ZapLog.Error(args...) }
func Errorw(msg string, kv ...interface{})        { ZapLog.Errorw(msg, kv...) }
func Errorf(template string, args ...interface{}) { ZapLog.Errorf(template, args...) }

// Panic / Fatal
func Panic(args ...interface{})                   { ZapLog.Panic(args...) }
func Panicf(template string, args ...interface{}) { ZapLog.Panicf(template, args...) }
func Fatal(args ...interface{})                   { ZapLog.Fatal(args...) }
func Fatalf(template string, args ...interface{}) { ZapLog.Fatalf(template, args...) }
