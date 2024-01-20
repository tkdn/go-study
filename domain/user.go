package domain

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db     *sqlx.DB
	tables struct {
		users  *goqu.SelectDataset
		musers *goqu.InsertDataset
	}
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
	r := &userRepo{}
	r.db = db
	r.tables.users = goqu.Dialect("postgres").From("users").Prepared(true)
	r.tables.musers = goqu.Dialect("postgres").Insert("users").Prepared(true)
	return r
}

func (r *userRepo) GetById(id int) (*User, error) {
	q, args, err := r.tables.users.Select("id", "name", "age").Where(goqu.C("id").Eq(id)).ToSQL()
	if err != nil {
		return nil, err
	}
	var user User
	if err := r.db.Get(&user, q, args...); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetList() ([]*User, error) {
	q, args, err := r.tables.users.Select("id", "name", "age").ToSQL()
	if err != nil {
		return nil, err
	}
	var users []*User
	if err := r.db.Select(&users, q, args...); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) Insert(name string, age int) (*User, error) {
	q, args, err := r.tables.musers.Rows(
		goqu.Record{"name": name, "age": age},
	).Returning("id", "name", "age").ToSQL()
	if err != nil {
		return nil, err
	}
	var user User
	if err := r.db.Get(&user, q, args...); err != nil {
		return nil, err
	}
	return &user, nil
}
