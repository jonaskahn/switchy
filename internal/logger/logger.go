package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.SugaredLogger
	once         sync.Once
)

// Init initializes the global logger with daily log rotation.
func Init() {
	once.Do(func() {
		logDir := filepath.Join(".", "logs")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("Failed to create logs directory: %v\n", err)
			globalLogger = zap.NewNop().Sugar()
			return
		}

		logFile := filepath.Join(logDir, "switchy.log")

		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100, // MB
			MaxBackups: 30,
			MaxAge:     30, // days
			Compress:   true,
			LocalTime:  true,
		}

		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

		consoleWriter := zapcore.AddSync(os.Stdout)
		fileWriter := zapcore.AddSync(lumberjackLogger)

		coreConfig := zap.NewProductionEncoderConfig()
		coreConfig.TimeKey = "timestamp"
		coreConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(coreConfig),
			fileWriter,
			zap.DebugLevel,
		)

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			consoleWriter,
			zap.InfoLevel,
		)

		core := zapcore.NewTee(fileCore, consoleCore)

		globalLogger = zap.New(core).Sugar()
	})
}

func GetLogger() *zap.SugaredLogger {
	if globalLogger == nil {
		Init()
	}
	return globalLogger
}

func Debug(args ...any) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...any) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...any) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...any) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...any) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...any) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...any) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...any) {
	GetLogger().Errorf(format, args...)
}

// Sync flushes any buffered log entries.
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}
