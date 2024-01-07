package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tkdn/go-study/log"
)

type userDB struct {
	db *sqlx.DB
}

type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
}

type UserRepo interface {
	GetById(id int) (*User, error)
}

func NewUserDB(db *sqlx.DB) UserRepo {
	return &userDB{db}
}

func (u *userDB) GetById(id int) (*User, error) {
	var user User
	err := u.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ConnectDB() *sqlx.DB {
	p, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dsn := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("DB_HOST"),
		p,
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Logger.Error(err.Error())
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Logger.Error(err.Error())
		panic(err)
	}

	return db
}
