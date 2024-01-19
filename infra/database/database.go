package database

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tkdn/go-study/log"
)

func ConnectDB() *sqlx.DB {
	db, err := sqlx.Open("pgx", GetDsn())
	if err != nil {
		log.Logger.Error(err.Error())
		panic(err)
	}
	if err := db.Ping(); err != nil {
		log.Logger.Error(err.Error())
		panic(err)
	}
	return db
}

func GetDsn() string {
	p, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		p,
		os.Getenv("POSTGRES_DB"),
	)
	return dsn
}
