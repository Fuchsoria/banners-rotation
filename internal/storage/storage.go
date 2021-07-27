package storage

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
