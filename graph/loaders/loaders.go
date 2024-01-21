package loaders

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/tkdn/go-study/domain"
)

type loaders struct {
	postLoader *dataloader.Loader[int, []*domain.Post]
}

const componentName = "github.com/tkdn/go-study/graph/loaders"

var (
	ctxKey          = &struct{ name string }{"loaders"}
	ErrPostNotFound = errors.New("post not found")
)

type loaderDeps struct {
	postRepo *domain.PostRepository
}

type Option func(*loaderDeps)

func WithPostRepo(r *domain.PostRepository) Option {
	return func(ld *loaderDeps) { ld.postRepo = r }
}

func New(opts ...Option) *loaders {
	var deps loaderDeps
	for _, apply := range opts {
		apply(&deps)
	}
	l := &loaders{
		postLoader: dataloader.NewBatchedLoader[int, []*domain.Post](
			getPostsByUserIDs(deps.postRepo),
			dataloader.WithClearCacheOnBatch[int, []*domain.Post](),
		),
	}
	return l
}

// impl func reterned actual loader result
func getPostsByUserIDs(repo *domain.PostRepository) dataloader.BatchFunc[int, []*domain.Post] {
	return func(ctx context.Context, userIDs []int) []*dataloader.Result[[]*domain.Post] {
		posts, err := repo.GetListByUserIDs(ctx, userIDs)
		if err != nil {
			posts = make(map[int][]*domain.Post)
		}
		result := make([]*dataloader.Result[[]*domain.Post], len(userIDs))
		for i, userID := range userIDs {
			res := new(dataloader.Result[[]*domain.Post])
			userPosts, ok := posts[userID]
			if ok {
				res.Data = userPosts
			} else {
				res.Error = ErrPostNotFound
			}
			result[i] = res
		}
		return result
	}
}

// retrun thunk for resolver
func GetPostsByUserID(ctx context.Context, userID int) ([]*domain.Post, error) {
	loaders := ctx.Value(ctxKey).(*loaders)
	thunk := loaders.postLoader.Load(ctx, userID)
	return thunk()
}

// impl as gglgen graphql.HandlerExtension and graphql.OperationInterceptor
var (
	_ graphql.HandlerExtension     = (*loaders)(nil)
	_ graphql.OperationInterceptor = (*loaders)(nil)
)

func (*loaders) ExtensionName() string { return componentName }

func (*loaders) Validate(graphql.ExecutableSchema) error { return nil }

func (r *loaders) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	return next(context.WithValue(ctx, ctxKey, r))
}
