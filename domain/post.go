package domain

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

type Post struct {
	ID        int       `db:"id" json:"id"`
	Text      string    `db:"text" json:"text"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
