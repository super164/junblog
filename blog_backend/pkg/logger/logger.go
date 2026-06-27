package logger

import (
	"os"

	"blog_backend/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

// Init 初始化日志
func Init(cfg *config.LogConfig) error {
	writerSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.FileName,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var encoder zapcore.Encoder
	if config.Get().App.Mode == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writerSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)

	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	return logger
}

// Sync 同步日志
func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
