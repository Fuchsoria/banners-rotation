package scripts

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	HTTPHost    = os.Getenv("TESTS_HTTP_HOST")
	PostgresDSN = os.Getenv("TESTS_POSTGRES_DSN")
)

func init() {
	if HTTPHost == "" {
		HTTPHost = "http://0.0.0.0:5555"
	}

	if PostgresDSN == "" {
		PostgresDSN = "host=0.0.0.0 port=5432 user=postgres password=example dbname=calendar sslmode=disable"
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

const delay = 5 * time.Second

func TestHTTP(t *testing.T) {
	log.Printf("wait %s for table creation...", delay)
	time.Sleep(delay)

	_, err := sqlx.ConnectContext(context.Background(), "postgres", PostgresDSN)
	if err != nil {
		panicOnErr(err)
	}

	t.Run("test add banner", func(t *testing.T) {
	})
}
