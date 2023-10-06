package storage

import (
	"context"
	"errors"
	"fmt"

	// third party
	"gorm.io/gorm"

	// external
	"github.com/taraslis453/solid-software-test/pkg/postgresql"

	// internal
	"github.com/taraslis453/solid-software-test/internal/entity"
	"github.com/taraslis453/solid-software-test/internal/service"
)

var _ service.UserStorage = (*userStorage)(nil)

type userStorage struct {
	*postgresql.PostgreSQLGorm
}

func NewUserStorage(postgresql *postgresql.PostgreSQLGorm) *userStorage {
	return &userStorage{postgresql}
}

func (r *userStorage) GetUser(ctx context.Context, filter service.GetUserFilter) (*entity.User, error) {
	stmt := r.DB
	if filter.EmailAddress != nil {
		stmt = stmt.Where(entity.User{EmailAddress: *filter.EmailAddress})
	}
	if filter.ID != nil {
		stmt = stmt.Where(entity.User{ID: *filter.ID})
	}

	var user entity.User
	err := stmt.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userStorage) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.DB.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *userStorage) UpdateUser(ctx context.Context, id string, user *entity.User) (*entity.User, error) {
	err := r.DB.Model(&entity.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *userStorage) DeleteUser(ctx context.Context, id string) error {
	err := r.DB.Delete(&entity.User{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
