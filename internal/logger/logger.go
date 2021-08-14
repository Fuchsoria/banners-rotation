package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	instance *zap.Logger
}

func New(level string, file string) *Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   file,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.TimeKey = "time"

	aLevel := zap.NewAtomicLevel()

	switch level {
	case "debug":
		aLevel.SetLevel(zapcore.DebugLevel)
	case "panic":
		aLevel.SetLevel(zapcore.PanicLevel)
	case "error":
		aLevel.SetLevel(zapcore.ErrorLevel)
	case "info":
		aLevel.SetLevel(zapcore.InfoLevel)
	case "fatal":
		aLevel.SetLevel(zapcore.FatalLevel)
	case "warn":
		aLevel.SetLevel(zapcore.WarnLevel)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		w,
		aLevel,
	)

	return &Logger{zap.New(core)}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.instance.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.instance.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.instance.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.instance.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) GetInstance() *zap.Logger {
	return l.instance
}
