package database

import (
	"github.com/jmoiron/sqlx"
	"github.com/tkdn/go-study/log"
)

type userRepo struct {
	db *sqlx.DB
}

// API レスポンスのモデルも兼ねているが
// Presenterのモデルとして切り出してもいいだろう
type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Age  int    `db:"age" json:"age"`
}

type UserRepository interface {
	GetById(id int) (*User, error)
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{db}
}

func (u *userRepo) GetById(id int) (*User, error) {
	var user User
	stmt, err := u.db.Preparex(`SELECT * FROM users WHERE id = $1`)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	if err := stmt.Get(&user, id); err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	return &user, nil
}
