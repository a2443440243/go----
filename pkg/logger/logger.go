package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// logrusLogger logrus实现
type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// New 创建新的日志实例
func New(level string) Logger {
	logger := logrus.New()

	// 设置日志级别
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 设置JSON格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置输出
	logger.SetOutput(os.Stdout)

	return &logrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug 调试日志
func (l *logrusLogger) Debug(msg string) {
	l.entry.Debug(msg)
}

// Info 信息日志
func (l *logrusLogger) Info(msg string) {
	l.entry.Info(msg)
}

// Warn 警告日志
func (l *logrusLogger) Warn(msg string) {
	l.entry.Warn(msg)
}

// Error 错误日志
func (l *logrusLogger) Error(msg string) {
	l.entry.Error(msg)
}

// Fatal 致命错误日志
func (l *logrusLogger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

// WithField 添加单个字段
func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields 添加多个字段
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}