package domain_test

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tkdn/go-study/infra"
	"github.com/tkdn/go-study/log"
)

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	db := infra.ConnectDB()
	testDB = db
	if _, err := db.Exec(`TRUNCATE TABLE users, posts`); err != nil {
		log.Logger.Error("failed to trucate tables")
	}
	os.Exit(m.Run())
}
