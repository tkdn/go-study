package domain

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
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

func NewPostRepository(db *sqlx.DB) *PostRepository {
	r := &PostRepository{}
	r.db = db
	r.tables.posts = goqu.Dialect("postgres").From("posts").Prepared(true)
	r.tables.mposts = goqu.Dialect("postgres").Insert("posts").Prepared(true)
	return r
}

func (r *PostRepository) GetListByUserIDs(ctx context.Context, userIDs []int) (map[int][]*Post, error) {
	q, args, err := r.tables.posts.Select("id", "user_id", "text", "created_at").Where(goqu.C("user_id").In(userIDs)).ToSQL()
	if err != nil {
		return nil, err
	}
	var posts []*Post
	if err := r.db.SelectContext(ctx, &posts, q, args...); err != nil {
		return nil, err
	}
	groupByUserIds := make(map[int][]*Post, 0)
	for _, post := range posts {
		if mapV, ok := groupByUserIds[post.UserID]; ok {
			groupByUserIds[post.UserID] = append(mapV, post)
			continue
		}
		groupByUserIds[post.UserID] = []*Post{post}
	}
	return groupByUserIds, err
}

func (r *PostRepository) GetListByUserID(ctx context.Context, id int) ([]*Post, error) {
	q, args, err := r.tables.posts.Select("id", "text", "created_at").Where(goqu.C("user_id").Eq(id)).ToSQL()
	if err != nil {
		return nil, err
	}
	var posts []*Post
	if err := r.db.SelectContext(ctx, &posts, q, args...); err != nil {
		return nil, err
	}
	return posts, err
}

func (r *PostRepository) Insert(ctx context.Context, userId int, text string) (*Post, error) {
	q, args, err := r.tables.mposts.Rows(
		goqu.Record{"user_id": userId, "text": text},
	).Returning("id", "user_id", "text", "created_at").ToSQL()
	if err != nil {
		return nil, err
	}
	var post Post
	if err := r.db.GetContext(ctx, &post, q, args...); err != nil {
		return nil, err
	}
	return &post, nil
}
