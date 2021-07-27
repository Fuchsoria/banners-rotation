package app

import (
	"github.com/Fuchsoria/banners-rotation/internal/storage"
	"go.uber.org/zap"
)

type App struct {
	logger  Logger
	storage Storage
	bandit  Bandit
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
	AddClickEvent(bannerID string, slotID string, socialDemoID string) error
	AddViewEvent(bannerID string, slotID string, socialDemoID string) error
	GetNotViewedBanners(slotID string) ([]storage.NotViewedItem, error)
	GetBannersClicks(slotID string) ([]storage.ClickItem, error)
	GetBannersViews(slotID string) ([]storage.ViewItem, error)
	GetBannersInSlot(slotID string) ([]storage.BannerRotationItem, error)
	CreateBanner(ID string, description string) (string, error)
	CreateSlot(ID string, description string) (string, error)
	CreateSocialDemo(ID string, description string) (string, error)
}

type Bandit interface {
	Use(items []string, clicks map[string]int, views map[string]int) string
}

func New(logger Logger, storage Storage, bandit Bandit) *App {
	return &App{logger, storage, bandit}
}

func (a *App) GetLogger() Logger {
	return a.logger
}

func (a *App) GetStorage() Storage {
	return a.storage
}

func (a *App) AddBannerRotation(bannerID string, slotID string) error {
	return a.storage.AddBannerRotation(bannerID, slotID)
}
func (a *App) RemoveBannerRotation(bannerID string, slotID string) error {
	return a.storage.RemoveBannerRotation(bannerID, slotID)
}
func (a *App) AddClickEvent(bannerID string, slotID string, socialDemoID string) error {
	return a.storage.AddClickEvent(bannerID, slotID, socialDemoID)
}
func (a *App) AddViewEvent(bannerID string, slotID string, socialDemoID string) error {
	return a.storage.AddViewEvent(bannerID, slotID, socialDemoID)
}
func (a *App) MapDataFromDB(
	bannersInSlot []storage.BannerRotationItem,
	bannersClicks []storage.ClickItem,
	bannersViews []storage.ViewItem) (
	banners []string,
	mappedBannersClicks map[string]int,
	mappedBannersViews map[string]int,
) {
	mappedBannersClicks = make(map[string]int)
	mappedBannersViews = make(map[string]int)

	for _, banner := range bannersInSlot {
		banners = append(banners, banner.BannerID)
	}

	for _, click := range bannersClicks {
		mappedBannersClicks[click.BannerID]++
	}

	for _, view := range bannersViews {
		mappedBannersViews[view.BannerID]++
	}

	return banners, mappedBannersClicks, mappedBannersClicks
}
func (a *App) GetBanner(slotID string, socialDemoID string) (string, error) {
	notViewedBanners, err := a.storage.GetNotViewedBanners(slotID)
	if err != nil {
		return "", err
	}

	if len(notViewedBanners) > 0 {
		bannerID := notViewedBanners[0].BannerID

		err := a.AddViewEvent(bannerID, slotID, socialDemoID)
		if err != nil {
			return "", err
		}

		return bannerID, nil
	}

	bannersInSlot, err := a.storage.GetBannersInSlot(slotID)
	if err != nil {
		return "", err
	}

	bannersClicks, err := a.storage.GetBannersClicks(slotID)
	if err != nil {
		return "", err
	}

	bannersViews, err := a.storage.GetBannersViews(slotID)
	if err != nil {
		return "", err
	}

	banners, mappedBannersClicks, mappedBannersViews := a.MapDataFromDB(bannersInSlot, bannersClicks, bannersViews)
	bannerID := a.bandit.Use(banners, mappedBannersClicks, mappedBannersViews)

	err = a.AddViewEvent(bannerID, slotID, socialDemoID)
	if err != nil {
		return "", err
	}

	return bannerID, nil
}
func (a *App) CreateBanner(id string, description string) (string, error) {
	return a.storage.CreateBanner(id, description)
}
func (a *App) CreateSlot(id string, description string) (string, error) {
	return a.storage.CreateSlot(id, description)
}
func (a *App) CreateSocialDemo(id string, description string) (string, error) {
	return a.storage.CreateSocialDemo(id, description)
}
