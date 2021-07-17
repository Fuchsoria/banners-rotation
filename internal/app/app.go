package app

import (
	"go.uber.org/zap"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	GetInstance() *zap.Logger
}

type Storage interface {
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) GetLogger() Logger {
	return a.logger
}
