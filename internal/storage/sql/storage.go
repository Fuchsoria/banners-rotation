package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	simpleproducer "github.com/Fuchsoria/banners-rotation/internal/amqp/producer"
	"github.com/Fuchsoria/banners-rotation/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Producer interface {
	Publish(message simpleproducer.AMQPMessage) error
}

type Storage struct {
	db       *sqlx.DB
	producer Producer
}

var ErrBannersWereRemoved = errors.New("banners were not removed from rotation")

func New(ctx context.Context, connectionString string, producer Producer) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot open db, %w", err)
	}

	return &Storage{db, producer}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("cannot connect to db, %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddBannerRotation(bannerID string, slotID string) error {
	_, err := s.db.Exec("INSERT INTO banners_rotation (slot_id,banner_id) VALUES ($1,$2)", slotID, bannerID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBannerRotation(bannerID string, slotID string) error {
	result, err := s.db.Exec("DELETE FROM banners_rotation WHERE slot_id=$1 AND banner_id=$2", slotID, bannerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBannersWereRemoved
	}

	return nil
}

func (s *Storage) AddClickEvent(bannerID string, slotID string, socialDemoID string) error {
	date := time.Now().String()
	_, err := s.db.Exec("INSERT INTO clicks (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
	if err != nil {
		return err
	}

	err = s.producer.Publish(simpleproducer.AMQPMessage{Type: "click", SlotID: slotID, BannerID: bannerID, SocialDemoID: socialDemoID, Date: date})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddViewEvent(bannerID string, slotID string, socialDemoID string) error {
	date := time.Now().String()
	_, err := s.db.Exec("INSERT INTO views (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
	if err != nil {
		return err
	}

	err = s.producer.Publish(simpleproducer.AMQPMessage{Type: "view", SlotID: slotID, BannerID: bannerID, SocialDemoID: socialDemoID, Date: date})
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetNotViewedBanners(slotID string) (notViewedBanners []storage.NotViewedItem, err error) {
	err = s.db.Select(&notViewedBanners, "SELECT slot_id,banner_id FROM banners_rotation WHERE slot_id=$1 EXCEPT SELECT slot_id,banner_id FROM views", slotID)
	if err != nil {
		return nil, err
	}

	return notViewedBanners, nil
}

func (s *Storage) GetBannersInSlot(slotID string) (bannersInSlot []storage.BannerRotationItem, err error) {
	err = s.db.Select(&bannersInSlot, "SELECT * FROM banners_rotation WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, err
	}

	return bannersInSlot, nil
}

func (s *Storage) GetBannersClicks(slotID string) (bannersClicks []storage.ClickItem, err error) {
	err = s.db.Select(&bannersClicks, "SELECT * FROM clicks WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, err
	}

	return bannersClicks, nil
}

func (s *Storage) GetBannersViews(slotID string) (bannersViews []storage.ViewItem, err error) {
	err = s.db.Select(&bannersViews, "SELECT * FROM views WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, err
	}

	return bannersViews, nil
}

func (s *Storage) CreateBanner(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO banners (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) CreateSlot(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO slots (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) CreateSocialDemo(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO social_demos (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", err
	}

	return id, nil
}
