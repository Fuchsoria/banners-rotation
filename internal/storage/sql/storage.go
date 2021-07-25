package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Bandit interface {
	Use(items []string, clicks map[string]int, views map[string]int) string
}

type Storage struct {
	db     *sqlx.DB
	bandit Bandit
}

type BannerRotationItem struct {
	SlotID   string `db:"slot_id"`
	BannerID string `db:"banner_id"`
}

type SessionClickItem struct {
	SlotID       string `db:"slot_id"`
	BannerID     string `db:"banner_id"`
	SocialDemoID string `db:"social_demo_id"`
	Date         string `db:"date"`
}

type SessionViewItem struct {
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

func New(ctx context.Context, connectionString string, bandit Bandit) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot open db, %w", err)
	}

	return &Storage{db, bandit}, nil
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

func (s *Storage) ClearSessionClicks() error {
	_, err := s.db.Exec("TRUNCATE TABLE session_clicks")
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ClearSessionViews() error {
	_, err := s.db.Exec("TRUNCATE TABLE session_views")
	if err != nil {
		return err
	}

	return nil
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

func (s *Storage) AddSessionClickEvent(bannerID string, slotID string, socialDemoID string) error {
	_, err := s.db.Exec("INSERT INTO session_clicks (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, time.Now().String())
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddSessionViewEvent(bannerID string, slotID string, socialDemoID string) error {
	_, err := s.db.Exec("INSERT INTO session_views (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, time.Now().String())
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) MapDataFromDB(
	bannersInSlot []BannerRotationItem,
	sessionBannersClicks []SessionClickItem,
	sessionBannersViews []SessionViewItem) (
	banners []string,
	bannersClicks map[string]int,
	bannersViews map[string]int,
) {
	bannersClicks = make(map[string]int)
	bannersViews = make(map[string]int)

	for _, banner := range bannersInSlot {
		banners = append(banners, banner.BannerID)
	}

	for _, click := range sessionBannersClicks {
		bannersClicks[click.BannerID]++
	}

	for _, view := range sessionBannersViews {
		bannersViews[view.BannerID]++
	}

	return banners, bannersClicks, bannersViews
}

func (s *Storage) GetBanner(slotID string, socialDemoID string) (string, error) {
	var bannersInSlot []BannerRotationItem
	var sessionBannersClicks []SessionClickItem
	var sessionBannersViews []SessionViewItem
	var notViewedBanners []NotViewedItem

	err := s.db.Select(&notViewedBanners, "SELECT slot_id,banner_id FROM banners_rotation WHERE slot_id=$1 EXCEPT SELECT slot_id,banner_id FROM session_views", slotID)
	if err != nil {
		return "", err
	}

	if len(notViewedBanners) > 0 {
		bannerID := notViewedBanners[0].BannerID

		err := s.AddSessionViewEvent(bannerID, slotID, socialDemoID)
		if err != nil {
			return "", err
		}

		return bannerID, nil
	}

	err = s.db.Select(&bannersInSlot, "SELECT * FROM banners_rotation WHERE slot_id=$1", slotID)
	if err != nil {
		return "", err
	}

	err = s.db.Select(&sessionBannersClicks, "SELECT * FROM session_clicks WHERE slot_id=$1", slotID)
	if err != nil {
		return "", err
	}

	err = s.db.Select(&sessionBannersViews, "SELECT * FROM session_views WHERE slot_id=$1", slotID)
	if err != nil {
		return "", err
	}

	banners, bannersClicks, bannersViews := s.MapDataFromDB(bannersInSlot, sessionBannersClicks, sessionBannersViews)
	bannerID := s.bandit.Use(banners, bannersClicks, bannersViews)

	err = s.AddSessionViewEvent(bannerID, slotID, socialDemoID)
	if err != nil {
		return "", err
	}

	return bannerID, nil
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
