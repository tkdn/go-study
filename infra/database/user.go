package database

import "github.com/jmoiron/sqlx"

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
	err := u.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
