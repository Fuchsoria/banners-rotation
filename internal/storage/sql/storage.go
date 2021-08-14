package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

type BannerRotationItem struct {
	SlotID   string `db:"slot_id"`
	BannerID string `db:"banner_id"`
}

type ClickItem struct {
	SlotID       string `db:"slot_id"`
	BannerID     string `db:"banner_id"`
	SocialDemoID string `db:"social_demo_id"`
	Date         string `db:"date"`
}

type ViewItem struct {
	SlotID       string `db:"slot_id"`
	BannerID     string `db:"banner_id"`
	SocialDemoID string `db:"social_demo_id"`
	Date         string `db:"date"`
}

type NotViewedItem struct {
	SlotID   string `db:"slot_id"`
	BannerID string `db:"banner_id"`
}

var ErrBannersWereRemoved = errors.New("banners were not removed from rotation")

func New(ctx context.Context, connectionString string) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot open db, %w", err)
	}

	return &Storage{db}, nil
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
		return fmt.Errorf("cannot insert banner to rotation, %w", err)
	}

	return nil
}

func (s *Storage) RemoveBannerRotation(bannerID string, slotID string) error {
	result, err := s.db.Exec("DELETE FROM banners_rotation WHERE slot_id=$1 AND banner_id=$2", slotID, bannerID)
	if err != nil {
		return fmt.Errorf("cannot delete banner from rotation, %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("cannot get affected rows count, %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rows are not affected on rotation delete, %w", ErrBannersWereRemoved)
	}

	return nil
}

func (s *Storage) AddClickEvent(bannerID string, slotID string, socialDemoID string, date string) error {
	_, err := s.db.Exec("INSERT INTO clicks (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
	if err != nil {
		return fmt.Errorf("cannot insert banner click, %w", err)
	}

	return nil
}

func (s *Storage) AddViewEvent(bannerID string, slotID string, socialDemoID string, date string) error {
	_, err := s.db.Exec("INSERT INTO views (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
	if err != nil {
		return fmt.Errorf("cannot insert banner view, %w", err)
	}

	return nil
}

func (s *Storage) GetNotViewedBanners(slotID string) (notViewedBanners []NotViewedItem, err error) {
	err = s.db.Select(&notViewedBanners, "SELECT slot_id,banner_id FROM banners_rotation WHERE slot_id=$1 EXCEPT SELECT slot_id,banner_id FROM views", slotID)
	if err != nil {
		return nil, fmt.Errorf("cannot get not viewed banners, %w", err)
	}

	return notViewedBanners, nil
}

func (s *Storage) GetBannersInSlot(slotID string) (bannersInSlot []BannerRotationItem, err error) {
	err = s.db.Select(&bannersInSlot, "SELECT * FROM banners_rotation WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, fmt.Errorf("cannot get banners from slot, %w", err)
	}

	return bannersInSlot, nil
}

func (s *Storage) GetBannersClicks(slotID string) (bannersClicks []ClickItem, err error) {
	err = s.db.Select(&bannersClicks, "SELECT * FROM clicks WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, fmt.Errorf("cannot get clicked banners, %w", err)
	}

	return bannersClicks, nil
}

func (s *Storage) GetBannersViews(slotID string) (bannersViews []ViewItem, err error) {
	err = s.db.Select(&bannersViews, "SELECT * FROM views WHERE slot_id=$1", slotID)
	if err != nil {
		return nil, fmt.Errorf("cannot get viewed banners, %w", err)
	}

	return bannersViews, nil
}

func (s *Storage) CreateBanner(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO banners (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", fmt.Errorf("cannot insert banner, %w", err)
	}

	return id, nil
}

func (s *Storage) CreateSlot(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO slots (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", fmt.Errorf("cannot insert slot, %w", err)
	}

	return id, nil
}

func (s *Storage) CreateSocialDemo(id string, description string) (string, error) {
	_, err := s.db.Exec("INSERT INTO social_demos (id,description) VALUES ($1,$2)", id, description)
	if err != nil {
		return "", fmt.Errorf("cannot insert social demo, %w", err)
	}

	return id, nil
}
