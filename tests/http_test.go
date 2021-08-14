package scripts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	simpleproducer "github.com/Fuchsoria/banners-rotation/internal/amqp/producer"
	sqlstorage "github.com/Fuchsoria/banners-rotation/internal/storage/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

type CreateBody struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type AddBannerBody struct {
	BannerID string `json:"banner_id"`
	SlotID   string `json:"slot_id"`
}

type RemoveBannerBody struct {
	BannerID string `json:"banner_id"`
	SlotID   string `json:"slot_id"`
}

type AddBannerClickBody struct {
	BannerID     string `json:"banner_id"`
	SlotID       string `json:"slot_id"`
	SocialDemoID string `json:"social_demo_id"`
}

type GetBannerBody struct {
	SlotID       string `json:"slot_id"`
	SocialDemoID string `json:"social_demo_id"`
}

type IDResponse struct {
	ID string `json:"id"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ItemDB struct {
	ID          string `db:"id"`
	Description string `db:"description"`
}

type RotationDB struct {
	BannerID string `db:"banner_id"`
	SlotID   string `db:"slot_id"`
}

type ClickDB struct {
	BannerID     string `db:"banner_id"`
	SlotID       string `db:"slot_id"`
	SocialDemoID string `db:"social_demo_id"`
	Date         string `db:"date"`
}

type ViewDB struct {
	BannerID     string `db:"banner_id"`
	SlotID       string `db:"slot_id"`
	SocialDemoID string `db:"social_demo_id"`
	Date         string `db:"date"`
}

var (
	HTTPHost    = os.Getenv("TESTS_HTTP_HOST")
	PostgresDSN = os.Getenv("TESTS_POSTGRES_DSN")
	AmpqDSN     = os.Getenv("TESTS_AMQP_DSN")
)

func init() {
	if HTTPHost == "" {
		HTTPHost = "http://0.0.0.0:5555"
	}

	if PostgresDSN == "" {
		PostgresDSN = "amqp://guest:guest@rabbit_test:5672/"
	}

	if AmpqDSN == "" {
		PostgresDSN = "host=0.0.0.0 port=5432 user=postgres password=example dbname=banners-rotation_test sslmode=disable"
	}
}

func TestStorage(t *testing.T) {
	conn, err := amqp.Dial(AmpqDSN)
	require.NoError(t, err)

	db, err := sqlx.ConnectContext(context.Background(), "postgres", PostgresDSN)
	require.NoError(t, err)

	producer := simpleproducer.New("banners-rotation", conn)
	err = producer.Connect()
	require.NoError(t, err, "should be without errors")

	storage, err := sqlstorage.New(context.Background(), PostgresDSN)
	require.NoError(t, err, "should be without errors")

	t.Run("test banner create", func(t *testing.T) {
		id := uuid.NewString()

		_, err := storage.CreateBanner(id, "")
		require.NoError(t, err, "should be without errors")

		var banner ItemDB

		err = db.Get(&banner, "SELECT * FROM banners WHERE id=$1", id)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, id, banner.ID, "item should be created in db")
	})

	t.Run("test slot create", func(t *testing.T) {
		id := uuid.NewString()

		_, err := storage.CreateSlot(id, "")
		require.NoError(t, err, "should be without errors")

		var slot ItemDB

		err = db.Get(&slot, "SELECT * FROM slots WHERE id=$1", id)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, id, slot.ID, "item should be created in db")
	})

	t.Run("test social demo create", func(t *testing.T) {
		id := uuid.NewString()

		_, err := storage.CreateSocialDemo(id, "")
		require.NoError(t, err, "should be without errors")

		var socialDemo ItemDB

		err = db.Get(&socialDemo, "SELECT * FROM social_demos WHERE id=$1", id)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, id, socialDemo.ID, "item should be created in db")
	})

	t.Run("test add banner to rotation", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		err := storage.AddBannerRotation(bannerID, slotID)
		require.NoError(t, err, "should be without errors")

		var rotation RotationDB

		err = db.Get(&rotation, "SELECT * FROM banners_rotation WHERE banner_id=$1 AND slot_id=$2", bannerID, slotID)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, bannerID, rotation.BannerID, "item should be created in db")
		require.Equal(t, slotID, rotation.SlotID, "item should be created in db")
	})

	t.Run("test remove banner from rotation", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		_, err := db.Query("INSERT INTO banners_rotation (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
		require.NoError(t, err, "should be without errors")

		err = storage.RemoveBannerRotation(bannerID, slotID)
		require.NoError(t, err, "should be without errors")

		var rotation []RotationDB

		err = db.Select(&rotation, "SELECT * FROM banners_rotation WHERE banner_id=$1 AND slot_id=$2", bannerID, slotID)
		require.NoError(t, err, "should be without errors")
		require.Len(t, rotation, 0, "selected rotation should be empty")
	})

	t.Run("test add banner click", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()
		socialDemoID := uuid.NewString()

		err := storage.AddClickEvent(bannerID, slotID, socialDemoID, time.Now().String())
		require.NoError(t, err, "should be without errors")

		var click ClickDB

		err = db.Get(&click, "SELECT * FROM clicks WHERE slot_id=$1 AND banner_id=$2 AND social_demo_id=$3", slotID, bannerID, socialDemoID)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, bannerID, click.BannerID, "item should be created in db")
		require.Equal(t, slotID, click.SlotID, "item should be created in db")
		require.Equal(t, socialDemoID, click.SocialDemoID, "item should be created in db")
		require.NotEmpty(t, click.Date, "date should exist")
	})

	t.Run("test add banner views", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()
		socialDemoID := uuid.NewString()

		err := storage.AddViewEvent(bannerID, slotID, socialDemoID, time.Now().String())
		require.NoError(t, err, "should be without errors")

		var view ViewDB

		err = db.Get(&view, "SELECT * FROM views WHERE slot_id=$1 AND banner_id=$2 AND social_demo_id=$3", slotID, bannerID, socialDemoID)
		require.NoError(t, err, "should be without errors")
		require.Equal(t, bannerID, view.BannerID, "item should be created in db")
		require.Equal(t, slotID, view.SlotID, "item should be created in db")
		require.Equal(t, socialDemoID, view.SocialDemoID, "item should be created in db")
		require.NotEmpty(t, view.Date, "date should exist")
	})

	t.Run("test get banners clicks", func(t *testing.T) {
		slotID := uuid.NewString()
		bannerID := uuid.NewString()
		socialDemoID := uuid.NewString()
		date := time.Now().String()

		_, err := db.Exec("INSERT INTO clicks (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
		require.NoError(t, err, "should be without errors")

		clicks, err := storage.GetBannersClicks(slotID)
		require.NoError(t, err, "should be without errors")

		require.Len(t, clicks, 1, "slice should have 1 item")
		require.NotEmpty(t, clicks[0], "first item should exist")
		require.Equal(t, bannerID, clicks[0].BannerID, "bannerID should be same")
		require.Equal(t, slotID, clicks[0].SlotID, "slotID should be same")
		require.Equal(t, socialDemoID, clicks[0].SocialDemoID, "socialDemoID should be same")
		require.Equal(t, date, clicks[0].Date, "date should be same")
	})

	t.Run("test get banners views", func(t *testing.T) {
		slotID := uuid.NewString()
		bannerID := uuid.NewString()
		socialDemoID := uuid.NewString()
		date := time.Now().String()

		_, err := db.Exec("INSERT INTO views (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, socialDemoID, date)
		require.NoError(t, err, "should be without errors")

		views, err := storage.GetBannersViews(slotID)
		require.NoError(t, err, "should be without errors")

		require.Len(t, views, 1, "slice should have 1 item")
		require.NotEmpty(t, views[0], "first item should exist")
		require.Equal(t, bannerID, views[0].BannerID, "bannerID should be same")
		require.Equal(t, slotID, views[0].SlotID, "slotID should be same")
		require.Equal(t, socialDemoID, views[0].SocialDemoID, "socialDemoID should be same")
		require.Equal(t, date, views[0].Date, "date should be same")
	})

	t.Run("test get not viewed banners", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		_, err := db.Query("INSERT INTO banners_rotation (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
		require.NoError(t, err, "should be without errors")

		notViewedBanners, err := storage.GetNotViewedBanners(slotID)
		require.NoError(t, err, "should be without errors")

		require.Len(t, notViewedBanners, 1)
		require.NotEmpty(t, notViewedBanners[0])
		require.Equal(t, bannerID, notViewedBanners[0].BannerID, "bannerID should be same")
		require.Equal(t, slotID, notViewedBanners[0].SlotID, "slotID should be same")
	})

	t.Run("test get not viewed banners with 0 unviewed items", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		_, err := db.Query("INSERT INTO banners_rotation (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
		require.NoError(t, err, "should be without errors")

		_, err = db.Exec("INSERT INTO views (slot_id,banner_id,social_demo_id,date) VALUES ($1,$2,$3,$4)", slotID, bannerID, "", "")
		require.NoError(t, err, "should be without errors")

		notViewedBanners, err := storage.GetNotViewedBanners(slotID)
		require.NoError(t, err, "should be without errors")

		require.Len(t, notViewedBanners, 0, "should be empty array")
	})

	t.Run("test get banners in slot", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		_, err := db.Query("INSERT INTO banners_rotation (slot_id, banner_id) VALUES ($1, $2)", slotID, bannerID)
		require.NoError(t, err, "should be without errors")

		bannersInSlot, err := storage.GetBannersInSlot(slotID)
		require.NoError(t, err, "should be without errors")

		require.Len(t, bannersInSlot, 1)
		require.Equal(t, bannerID, bannersInSlot[0].BannerID, "bannerID should be same")
		require.Equal(t, slotID, bannersInSlot[0].SlotID, "slotID should be same")
	})
}

func TestHTTP(t *testing.T) {
	httpCreateBanner := HTTPHost + "/api/v1/admin/banners/create"
	httpCreateSlot := HTTPHost + "/api/v1/admin/slots/create"
	httpCreateSocialDemo := HTTPHost + "/api/v1/admin/social-demos/create"
	httpAddBanner := HTTPHost + "/api/v1/banners/add"
	httpRemoveBanner := HTTPHost + "/api/v1/banners/remove"
	httpAddBannerClick := HTTPHost + "/api/v1/banners/click"
	httpGetBanner := HTTPHost + "/api/v1/banners/get"

	t.Run("test banner create", func(t *testing.T) {
		id := uuid.NewString()

		jsonData, err := json.Marshal(CreateBody{ID: id, Description: ""})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpCreateBanner, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response IDResponse

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err, "should be without errors")

		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.Equal(t, id, response.ID, "response id should be equal")
	})

	t.Run("test slot create", func(t *testing.T) {
		id := uuid.NewString()

		jsonData, err := json.Marshal(CreateBody{ID: id, Description: ""})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpCreateSlot, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response IDResponse

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err, "should be without errors")

		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.Equal(t, id, response.ID, "response id should be equal")
	})

	t.Run("test social-demo create", func(t *testing.T) {
		id := uuid.NewString()

		jsonData, err := json.Marshal(CreateBody{ID: id, Description: ""})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpCreateSocialDemo, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response IDResponse

		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err, "should be without errors")

		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.Equal(t, id, response.ID, "response id should be equal")
	})

	t.Run("test add banner to rotation", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		jsonData, err := json.Marshal(AddBannerBody{BannerID: bannerID, SlotID: slotID})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpAddBanner, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response MessageResponse

		err = json.NewDecoder(resp.Body).Decode(&response)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.NotEmpty(t, response.Message, "response should exist")
	})

	t.Run("test remove banner from rotation", func(t *testing.T) {
		jsonData, err := json.Marshal(RemoveBannerBody{BannerID: "banner1", SlotID: "slot1"})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpRemoveBanner, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response MessageResponse

		err = json.NewDecoder(resp.Body).Decode(&response)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.NotEmpty(t, response.Message, "response should exist")
	})

	t.Run("test add banner click", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()
		socialDemoID := uuid.NewString()

		jsonData, err := json.Marshal(AddBannerClickBody{BannerID: bannerID, SlotID: slotID, SocialDemoID: socialDemoID})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpAddBannerClick, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		var response MessageResponse

		err = json.NewDecoder(resp.Body).Decode(&response)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.NotEmpty(t, response.Message, "response should exist")
	})

	t.Run("test get banner", func(t *testing.T) {
		socialDemoID := uuid.NewString()

		jsonData, err := json.Marshal(GetBannerBody{SlotID: "slot1", SocialDemoID: socialDemoID})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpGetBanner, "application/json",
			bytes.NewBuffer(jsonData))

		var response IDResponse

		json.NewDecoder(resp.Body).Decode(&response)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.NotEmpty(t, response.ID, "banner id should exist")
	})

	t.Run("test get banner from not existed slot", func(t *testing.T) {
		slotID := uuid.NewString()
		socialDemoID := uuid.NewString()

		jsonData, err := json.Marshal(GetBannerBody{SlotID: slotID, SocialDemoID: socialDemoID})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpGetBanner, "application/json",
			bytes.NewBuffer(jsonData))

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "response statuscode should be not found")
	})

	t.Run("test remove not existed banner from rotation", func(t *testing.T) {
		bannerID := uuid.NewString()
		slotID := uuid.NewString()

		jsonData, err := json.Marshal(RemoveBannerBody{BannerID: bannerID, SlotID: slotID})
		require.NoError(t, err, "should be without errors")

		resp, err := http.Post(httpRemoveBanner, "application/json",
			bytes.NewBuffer(jsonData))
		require.NoError(t, err, "should be without errors")

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusNotFound, resp.StatusCode, "response statuscode should be bad request")
	})

	t.Run("test empty body add banner", func(t *testing.T) {
		resp, err := http.Post(httpAddBanner, "application/json",
			nil)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "response statuscode should be bad request")
	})

	t.Run("test empty body remove banner", func(t *testing.T) {
		resp, err := http.Post(httpRemoveBanner, "application/json",
			nil)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "response statuscode should be bad request")
	})

	t.Run("test empty body add click", func(t *testing.T) {
		resp, err := http.Post(httpAddBannerClick, "application/json",
			nil)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "response statuscode should be bad request")
	})

	t.Run("test empty body get banner", func(t *testing.T) {
		resp, err := http.Post(httpGetBanner, "application/json",
			nil)

		require.NoError(t, err, "should be without errors")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode, "response statuscode should be bad request")
	})
}
