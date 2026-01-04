package logger

import (
	"os"

	"github.com/datpham2001/techblog/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func Initalize(cfg *config.Config) *Logger {
	l := &Logger{
		logger: logrus.New(),
	}

	l.init(cfg.Server.Env)
	return l
}

func (l *Logger) init(env string) {
	l.logger.SetOutput(os.Stdout)

	if env == "production" {
		l.logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		l.logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set log level based on environment
	switch env {
	case "production":
		l.logger.SetLevel(logrus.InfoLevel)
	default:
		l.logger.SetLevel(logrus.DebugLevel)
	}
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

func (l *Logger) Debug(args ...any) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Info(args ...any) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warn(args ...any) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Error(args ...any) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...any) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.logger.Panicf(format, args...)
}
