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
	GetList() ([]*User, error)
	Insert(name string, age int) (*User, error)
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

func (u *userRepo) GetList() ([]*User, error) {
	var users []*User
	rows, err := u.db.Queryx(`SELECT * FROM users`)
	if err != nil {
		log.Logger.Error(err.Error())
		return nil, err
	}
	for rows.Next() {
		var u User
		err = rows.StructScan(&u)
		if err != nil {
			log.Logger.Error(err.Error())
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (u *userRepo) Insert(name string, age int) (*User, error) {
	var user User
	stmt, err := u.db.Preparex(`INSERT INTO users(name, age) VALUES($1, $2) RETURNING id, name, age`)
	if err != nil {
		return nil, err
	}
	err = stmt.Get(&user, name, age)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
