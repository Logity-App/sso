package auth

import (
	"context"
	"github.com/Logity-App/sso/internal/infrastracture/repository"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	//"os/user"
)

//
//const (
//	UsersTable = `users`
//)

type Repository struct {
	client repository.Client
	//hashGenerator repository.HashGenerator
}

func NewUserRepository(client repository.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) CheckCredentials(ctx context.Context) error {
	//todo implement
	panic("imlement me!!")
}
