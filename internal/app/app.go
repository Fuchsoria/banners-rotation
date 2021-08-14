package app

import (
	"fmt"
	"time"

	simpleproducer "github.com/Fuchsoria/banners-rotation/internal/amqp/producer"
	sqlstorage "github.com/Fuchsoria/banners-rotation/internal/storage/sql"
	"go.uber.org/zap"
)

type App struct {
	logger   Logger
	storage  Storage
	bandit   Bandit
	producer Producer
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
	AddClickEvent(bannerID string, slotID string, socialDemoID string, date string) error
	AddViewEvent(bannerID string, slotID string, socialDemoID string, date string) error
	GetNotViewedBanners(slotID string) ([]sqlstorage.NotViewedItem, error)
	GetBannersClicks(slotID string) ([]sqlstorage.ClickItem, error)
	GetBannersViews(slotID string) ([]sqlstorage.ViewItem, error)
	GetBannersInSlot(slotID string) ([]sqlstorage.BannerRotationItem, error)
	CreateBanner(ID string, description string) (string, error)
	CreateSlot(ID string, description string) (string, error)
	CreateSocialDemo(ID string, description string) (string, error)
}

type Producer interface {
	Publish(message simpleproducer.AMQPMessage) error
}

type Bandit interface {
	Use(items []string, clicks map[string]int, views map[string]int) (string, error)
}

func New(logger Logger, storage Storage, bandit Bandit, producer Producer) *App {
	return &App{logger, storage, bandit, producer}
}

func (a *App) GetLogger() Logger {
	return a.logger
}

func (a *App) AddBannerRotation(bannerID string, slotID string) error {
	return a.storage.AddBannerRotation(bannerID, slotID)
}

func (a *App) RemoveBannerRotation(bannerID string, slotID string) error {
	return a.storage.RemoveBannerRotation(bannerID, slotID)
}

func (a *App) AddClickEvent(bannerID string, slotID string, socialDemoID string) error {
	date := time.Now().String()

	err := a.storage.AddClickEvent(bannerID, slotID, socialDemoID, date)
	if err != nil {
		return fmt.Errorf("cannot create banner click event, %w", err)
	}

	err = a.producer.Publish(simpleproducer.AMQPMessage{Type: "click", SlotID: slotID, BannerID: bannerID, SocialDemoID: socialDemoID, Date: date})
	if err != nil {
		return fmt.Errorf("cannot publish banner click, %w", err)
	}

	return nil
}

func (a *App) AddViewEvent(bannerID string, slotID string, socialDemoID string) error {
	date := time.Now().String()

	err := a.storage.AddViewEvent(bannerID, slotID, socialDemoID, date)
	if err != nil {
		return fmt.Errorf("cannot create banner view event, %w", err)
	}

	err = a.producer.Publish(simpleproducer.AMQPMessage{Type: "view", SlotID: slotID, BannerID: bannerID, SocialDemoID: socialDemoID, Date: date})
	if err != nil {
		return fmt.Errorf("cannot publish banner click, %w", err)
	}

	return nil
}

func (a *App) MapDataFromDB(
	bannersInSlot []sqlstorage.BannerRotationItem,
	bannersClicks []sqlstorage.ClickItem,
	bannersViews []sqlstorage.ViewItem) (
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

	return banners, mappedBannersClicks, mappedBannersViews
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
	bannerID, err := a.bandit.Use(banners, mappedBannersClicks, mappedBannersViews)
	if err != nil {
		return "", err
	}

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
