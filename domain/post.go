package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db     *sqlx.DB
	tables struct {
		posts  *goqu.SelectDataset
		mposts *goqu.InsertDataset
	}
}

type Post struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"userId"`
	Text      string    `db:"text" json:"text"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type PostRepositry interface {
	GetListByUserID(id int) ([]*Post, error)
	Insert(userID int, text string) (*Post, error)
}

func NewPostRepository(db *sqlx.DB) PostRepositry {
	r := &postRepo{}
	r.db = db
	r.tables.posts = goqu.Dialect("postgres").From("posts").Prepared(true)
	r.tables.mposts = goqu.Dialect("postgres").Insert("posts").Prepared(true)
	return r
}

func (r *postRepo) GetListByUserID(id int) ([]*Post, error) {
	q, args, err := r.tables.posts.Select("id", "text", "created_at").Where(goqu.C("user_id").Eq(id)).ToSQL()
	if err != nil {
		return nil, err
	}
	var posts []*Post
	if err := r.db.Select(&posts, q, args...); err != nil {
		return nil, err
	}
	return posts, err
}

func (r *postRepo) Insert(userId int, text string) (*Post, error) {
	q, args, err := r.tables.mposts.Rows(
		goqu.Record{"user_id": userId, "text": text},
	).Returning("id", "user_id", "text", "created_at").ToSQL()
	if err != nil {
		return nil, err
	}
	var post Post
	if err := r.db.Get(&post, q, args...); err != nil {
		return nil, err
	}
	return &post, nil
}
