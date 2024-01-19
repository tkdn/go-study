package graph

import (
	"github.com/tkdn/go-study/infra/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepo database.UserRepository
}
