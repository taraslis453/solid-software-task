package service

import (
	"context"

	"github.com/taraslis453/solid-software-test/internal/entity"
)

type Storages struct {
	User UserStorage
}

type UserStorage interface {
	GetUser(ctx context.Context, filter GetUserFilter) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, id string, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type GetUserFilter struct {
	ID           *string
	EmailAddress *string
}
