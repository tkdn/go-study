package infra

import (
	"fmt"
	"os"
	"strconv"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var (
	dbHost    = os.Getenv("DB_HOST")
	dbUser    = os.Getenv("POSTGRES_USER")
	dbPass    = os.Getenv("POSTGRES_PASSWORD")
	dbName    = os.Getenv("POSTGRES_DB")
	dbPort, _ = strconv.Atoi(os.Getenv("DB_PORT"))
)

func ConnectDB() (*sqlx.DB, error) {
	db, err := otelsql.Open("pgx", GetDsn(), otelsql.WithAttributes(newAttrs()...))
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(db, "pgx"), nil
}

func GetDsn() string {
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
	return dsn
}

func newAttrs() []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, 5)
	attrs = append(attrs,
		semconv.ServerAddress(dbHost),
		semconv.DBUser(dbUser),
		semconv.NetworkTransportTCP,
		semconv.ServerPort(dbPort),
		semconv.DBSystemPostgreSQL)
	return attrs
}
