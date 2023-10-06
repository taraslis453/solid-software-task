package service

import (
	"context"

	"github.com/taraslis453/solid-software-test/config"
	"github.com/taraslis453/solid-software-test/pkg/errs"
	"github.com/taraslis453/solid-software-test/pkg/logging"
	"github.com/taraslis453/solid-software-test/pkg/password"

	"github.com/taraslis453/solid-software-test/internal/entity"
)

type Services struct {
	User UserService
}

// serviceContext provides a shared context for all services
type serviceContext struct {
	storages Storages
	cfg      *config.Config
	logger   logging.Logger
}

// Options is used to parameterize service
type Options struct {
	Storages       Storages
	Config         *config.Config
	Logger         logging.Logger
	PasswordHasher password.Hasher
}

const (
	userNotFoundErrCode      = "user_not_found"
	userAlreadyExistsErrCode = "user_already_exists"

	invalidPasswordErrCode = "invalid_password"

	invalidTokenErrCode = "invalid_token"
	tokenExpiredErrCode = "token_expired"
)

type UserService interface {
	// RegisterUser is used to register a new user.
	RegisterUser(ctx context.Context, opt RegisterUserOptions) error
	// LoginUser is used to login a user.
	LoginUser(ctx context.Context, opt LoginUserOptions) (LoginUserOutput, error)
	// VerifyUserToken is used to verify the user by given token and return verified user entity.
	GetUser(ctx context.Context, opt GetUserOptions) (*entity.User, error)
	// UpdateUser is used to update a user.
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	// DeleteUser is used to delete a user.
	DeleteUser(ctx context.Context, id string) error
	// VerifyUserToken is used to verify the user by given token and return verified user entity.
	VerifyUserToken(ctx context.Context, token string) (*entity.User, error)
	// RefreshUserToken is used to verify refresh token and then generate a new pair of access and refresh tokens.
	RefreshUserToken(ctx context.Context, tokenStr string) (*UserTokenOutput, error)
	// GenerateUserToken is used to generate a pair of access and refresh tokens.
	GenerateUserToken(ctx context.Context, user *entity.User) (*UserTokenOutput, error)
}

var (
	ErrRegisterUserUserAlreadyExists = errs.New("user already exists", userAlreadyExistsErrCode)

	ErrLoginUserUserNotFound    = errs.New("user not found", userNotFoundErrCode)
	ErrLoginUserInvalidPassword = errs.New("invalid password", invalidPasswordErrCode)

	ErrGetUserUserNotFound = errs.New("user not found", userNotFoundErrCode)

	ErrUpdateUserUserNotFound = errs.New("user not found", userNotFoundErrCode)

	ErrDeleteUserUserNotFound = errs.New("user not found", userNotFoundErrCode)

	ErrVerifyUserTokenInvalidToken = errs.New("invalid authenticate token.", invalidTokenErrCode)
	ErrVerifyUserTokenUserNotFound = errs.New("user not found", userNotFoundErrCode)

	ErrRefreshUserTokenInvalidToken = errs.New("invalid refresh token", invalidTokenErrCode)
	ErrRefreshUserTokenUserNotFound = errs.New("user not found", userNotFoundErrCode)
)

type RegisterUserOptions struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Phone        string `json:"phone"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

type LoginUserOptions struct {
	EmailAddress string
	Password     string
}

type LoginUserOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

type VerifyTokenOptions struct {
	Token      string
	HMACSecret string
}

type GetUserOptions struct {
	ID           string
	EmailAddress string
}

type UserTokenOutput struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
