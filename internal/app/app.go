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
	AddBannerRotation(bannerID string, slotID string) error
	RemoveBannerRotation(bannerID string, slotID string) error
	AddSessionClickEvent(bannerID string, slotID string, socialDemoID string) error
	GetBanner(slotID string, socialDemoID string) (string, error)
	CreateBanner(ID string, description string) (string, error)
	CreateSlot(ID string, description string) (string, error)
	CreateSocialDemo(ID string, description string) (string, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) GetLogger() Logger {
	return a.logger
}

func (a *App) GetStorage() Storage {
	return a.storage
}
