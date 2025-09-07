package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/barimehdi77/cupid-api/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger initializes the global logger instance with enhanced readability
func InitLogger() error {
	// Get log level from environment (default: debug)
	logLevel := strings.ToLower(env.GetEnvString("LOG_LEVEL", "debug"))

	// Get environment (development or production)
	environment := strings.ToLower(env.GetEnvString("GO_ENV", "development"))

	var core zapcore.Core
	var err error

	if environment == "production" {
		// Production configuration: JSON output, optimized for performance
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}

		// Set log level
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))

		Logger, err = config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		// Development configuration: Enhanced human-readable output
		core = createDevelopmentCore(parseLogLevel(logLevel))
		Logger = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
			zap.Development(),
		)
	}

	if err != nil {
		return err
	}

	return nil
}

// createDevelopmentCore creates a highly readable console encoder for development
func createDevelopmentCore(level zapcore.Level) zapcore.Core {
	// Create a custom encoder config for maximum readability
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    customLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder,
	}

	// Create console encoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create writer syncer
	writeSyncer := zapcore.AddSync(os.Stdout)

	// Create atomic level
	atomicLevel := zap.NewAtomicLevelAt(level)

	// Return core
	return zapcore.NewCore(encoder, writeSyncer, atomicLevel)
}

// customLevelEncoder provides colored and padded level names for better readability
func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var levelStr string
	switch level {
	case zapcore.DebugLevel:
		levelStr = "\033[36m[DEBUG]\033[0m" // Cyan
	case zapcore.InfoLevel:
		levelStr = "\033[32m[INFO] \033[0m" // Green
	case zapcore.WarnLevel:
		levelStr = "\033[33m[WARN] \033[0m" // Yellow
	case zapcore.ErrorLevel:
		levelStr = "\033[31m[ERROR] \033[0m" // Red
	case zapcore.FatalLevel:
		levelStr = "\033[35m[FATAL] \033[0m" // Magenta
	default:
		levelStr = "[UNKNOWN]"
	}
	enc.AppendString(levelStr)
}

// customTimeEncoder provides a clean, readable timestamp format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05.000"))
}

// customCallerEncoder provides clean file:line information
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !caller.Defined {
		enc.AppendString("undefined")
		return
	}

	// Extract just the filename and line number for cleaner output
	file := caller.File
	if len(file) > 30 {
		// Truncate long paths, keep the last part
		parts := strings.Split(file, "/")
		if len(parts) > 2 {
			file = ".../" + strings.Join(parts[len(parts)-2:], "/")
		}
	}

	enc.AppendString(fmt.Sprintf("%s:%d", file, caller.Line))
}

// parseLogLevel converts string log level to zapcore.Level
func parseLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn", "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Helper functions for common logging operations

// Debug logs a debug message with optional fields
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug("🔍 "+msg, fields...)
}

// Info logs an info message with optional fields
func Info(msg string, fields ...zap.Field) {
	Logger.Info("ℹ️  "+msg, fields...)
}

// Warn logs a warning message with optional fields
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn("⚠️  "+msg, fields...)
}

// Error logs an error message with optional fields
func Error(msg string, fields ...zap.Field) {
	Logger.Error("❌ "+msg, fields...)
}

// Fatal logs a fatal message with optional fields and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal("💀 "+msg, fields...)
}

// With creates a child logger with the given fields
func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}

// Named creates a named logger
func Named(name string) *zap.Logger {
	return Logger.Named(name)
}

// Enhanced helper functions for better structured logging

// LogRequest logs HTTP request information in a structured way
func LogRequest(method, path string, statusCode int, duration time.Duration, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Duration("duration", duration),
	}

	allFields := append(baseFields, fields...)

	var icon string
	switch {
	case statusCode >= 500:
		icon = "🔥"
		Logger.Error(icon+" HTTP Request", allFields...)
	case statusCode >= 400:
		icon = "⚠️"
		Logger.Warn(icon+" HTTP Request", allFields...)
	case statusCode >= 300:
		icon = "🔄"
		Logger.Info(icon+" HTTP Request", allFields...)
	default:
		icon = "✅"
		Logger.Info(icon+" HTTP Request", allFields...)
	}
}

// LogStartup logs application startup information
func LogStartup(component string, fields ...zap.Field) {
	Logger.Info("🚀 "+component+" starting", fields...)
}

// LogShutdown logs application shutdown information
func LogShutdown(component string, fields ...zap.Field) {
	Logger.Info("🛑 "+component+" shutting down", fields...)
}

// LogError logs detailed error information
func LogError(operation string, err error, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("operation", operation),
		zap.Error(err),
	}
	allFields := append(baseFields, fields...)
	Logger.Error("❌ Operation failed", allFields...)
}

// LogSuccess logs successful operations
func LogSuccess(operation string, fields ...zap.Field) {
	Logger.Info("✅ "+operation+" completed successfully", fields...)
}

// LogProgress logs operation progress
func LogProgress(operation string, fields ...zap.Field) {
	Logger.Info("⏳ "+operation+" in progress", fields...)
}

// LogDatabase logs database operations
func LogDatabase(operation string, table string, duration time.Duration, fields ...zap.Field) {
	baseFields := []zap.Field{
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration),
	}
	allFields := append(baseFields, fields...)
	Logger.Debug("🗄️  Database operation", allFields...)
}
