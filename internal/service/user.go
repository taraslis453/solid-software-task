package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/taraslis453/solid-software-test/pkg/errs"
	"github.com/taraslis453/solid-software-test/pkg/password"
	"github.com/taraslis453/solid-software-test/pkg/token"

	"github.com/taraslis453/solid-software-test/internal/entity"
)

var _ UserService = (*userService)(nil)

type userService struct {
	serviceContext
	passwordHasher password.Hasher
}

func NewUserService(options Options) *userService {
	return &userService{
		serviceContext: serviceContext{
			storages: options.Storages,
			cfg:      options.Config,
			logger:   options.Logger.Named("userService"),
		},
		passwordHasher: options.PasswordHasher,
	}
}

func (s *userService) RegisterUser(ctx context.Context, opts RegisterUserOptions) error {
	logger := s.logger.
		Named("RegisterUser").
		WithContext(ctx).
		With("opts", opts)

	// Check if user with the same email already exists in storage and return error if exists
	user, err := s.storages.User.GetUser(ctx, GetUserFilter{
		EmailAddress: &opts.EmailAddress,
	})
	if err != nil {
		logger.Error("failed to get user through storage", "err", err)
		return fmt.Errorf("failed to get user through storage: %w", err)
	}
	if user != nil {
		logger.Info("user with such email already exists")
		return ErrRegisterUserUserAlreadyExists
	}

	hashedPassword, err := s.passwordHasher.GenerateHashFromPassword(opts.Password)
	if err != nil {
		logger.Error("failed to hash password", "err", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := s.storages.User.CreateUser(ctx, &entity.User{
		Name:         opts.Name,
		Surname:      opts.Surname,
		EmailAddress: opts.EmailAddress,
		Password:     hashedPassword,
	})
	if err != nil {
		logger.Error("failed to create user in storage", "err", err)
		return fmt.Errorf("failed to create user in storage: %w", err)
	}
	logger.Debug("created user", "createdUser", createdUser)

	logger.Info("registered user succesfully")
	return nil
}

func (s *userService) LoginUser(ctx context.Context, opts LoginUserOptions) (LoginUserOutput, error) {
	logger := s.logger.
		Named("LoginUser").
		WithContext(ctx).
		With("opts", opts)

	user, err := s.storages.User.GetUser(ctx, GetUserFilter{
		EmailAddress: &opts.EmailAddress,
	})
	if err != nil {
		logger.Error("failed to get user through storage", "err", err)
		return LoginUserOutput{}, fmt.Errorf("failed to get user through storage: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return LoginUserOutput{}, ErrLoginUserUserNotFound
	}
	logger.Debug("got user", "user", user)

	// Check if password is correct
	isPasswordCorrect, err := s.passwordHasher.CompareHashAndPassword(&password.CompareHashAndPasswordOptions{
		Hashed:   user.Password,
		Password: opts.Password,
	})
	if err != nil {
		logger.Error("failed to check password correctness", "err", err)
		return LoginUserOutput{}, fmt.Errorf("failed to check password correctness: %w", err)
	}
	if !isPasswordCorrect {
		logger.Info("invalid password")
		return LoginUserOutput{}, ErrLoginUserInvalidPassword
	}

	// Generate a pair of tokens
	tokens, err := s.GenerateUserToken(ctx, user)
	if err != nil {
		logger.Error("failed to generate tokens", "err", err)
		return LoginUserOutput{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	logger.Info("user logged in", "tokens", tokens)
	return LoginUserOutput{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserID:       user.ID,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, opt GetUserOptions) (*entity.User, error) {
	logger := s.logger.
		Named("GetUser").
		WithContext(ctx)

	user, err := s.storages.User.GetUser(ctx, GetUserFilter{
		ID:           &opt.ID,
		EmailAddress: &opt.EmailAddress,
	})
	if err != nil {
		logger.Error("failed to get user", "err", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return nil, ErrGetUserUserNotFound
	}
	logger = logger.With("user", user)
	logger.Debug("got user")

	logger.Info("successfully got user")
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, newUser *entity.User) (*entity.User, error) {
	logger := s.logger.
		Named("UpdateUser").
		With("newUser", newUser).
		WithContext(ctx)

	user, err := s.storages.User.GetUser(ctx, GetUserFilter{
		ID: &newUser.ID,
	})
	if err != nil {
		logger.Error("failed to get user", "err", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return nil, ErrUpdateUserUserNotFound
	}
	logger.Debug("got user")

	updatedUser, err := s.storages.User.UpdateUser(ctx, user.ID, newUser)
	if err != nil {
		logger.Error("failed to update user", "err", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	logger = logger.With("updatedUser", updatedUser)

	logger.Info("successfully updated user")
	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	logger := s.logger.
		Named("DeleteUser").
		With("id", id).
		WithContext(ctx)

	user, err := s.storages.User.GetUser(ctx, GetUserFilter{
		ID: &id,
	})
	if err != nil {
		logger.Error("failed to get user", "err", err)
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return ErrDeleteUserUserNotFound
	}
	logger.Debug("got user")

	err = s.storages.User.DeleteUser(ctx, id)
	if err != nil {
		logger.Error("failed to delete user", "err", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logger.Info("successfully deleted user")
	return nil
}

func (s *userService) RefreshUserToken(ctx context.Context, refreshToken string) (*UserTokenOutput, error) {
	logger := s.logger.
		Named("RefreshUserToken").
		WithContext(ctx).
		With("refreshToken", refreshToken)

	user, err := s.VerifyUserToken(ctx, refreshToken)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
		}
		switch err {
		case ErrVerifyUserTokenInvalidToken:
			return nil, ErrRefreshUserTokenInvalidToken
		case ErrVerifyUserTokenUserNotFound:
			return nil, ErrRefreshUserTokenUserNotFound
		default:
			logger.Error("failed to verify token", "err", err)
			return nil, fmt.Errorf("failed to verify token: %w", err)
		}
	}
	logger.Debug("got user", "user", user)

	tokens, err := s.GenerateUserToken(ctx, user)
	if err != nil {
		logger.Error("failed to generate token", "err", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("successfully refreshed tokens", "tokens", tokens)
	return tokens, nil
}

func (s *userService) VerifyUserToken(ctx context.Context, tokenStr string) (*entity.User, error) {
	logger := s.logger.
		Named("VerifyUserToken").
		WithContext(ctx).
		With("token", tokenStr)

	claims, err := token.VerifyJWTToken(tokenStr, s.cfg.Auth.TokenSecretKey)
	if err != nil {
		logger.Info("invalid token", "err", err)
		return nil, ErrVerifyUserTokenInvalidToken
	}
	var claimsData token.UserDataClaims
	if err := mapstructure.Decode(claims.GetPayload(), &claimsData); err != nil {
		return nil, ErrVerifyUserTokenInvalidToken
	}

	// Get user from storage by token &claims.UserID
	user, err := s.storages.User.GetUser(ctx, GetUserFilter{ID: &claimsData.UserID})
	if err != nil {
		logger.Error("failed to get user", "err", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return nil, ErrVerifyUserTokenUserNotFound
	}

	logger.Info("verified token", "user", user)
	return user, nil
}

func (ts *userService) GenerateUserToken(ctx context.Context, user *entity.User) (*UserTokenOutput, error) {
	logger := ts.logger.
		Named("GenerateUserToken").
		WithContext(ctx).
		With("user", user)

	// Create new Access token
	t := time.Now()
	accessToken, err := token.SignJWTToken(
		&token.UniversalClaims{
			Iss:   ts.cfg.Auth.TokenIssuer,
			ExpAt: t.Add(ts.cfg.Auth.AccessTokenLifetime),
			NbfAt: t,
			IssAt: t,
			Payload: token.UserDataClaims{
				UserID: user.ID,
			},
		},
		ts.cfg.Auth.TokenSecretKey,
	)
	if err != nil {
		logger.Error("failed to generate access token", "err", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	logger.Debug("access token generated", "accessToken", accessToken)

	// Create new Refresh token
	t = time.Now()
	refreshToken, err := token.SignJWTToken(
		&token.UniversalClaims{
			Iss:   ts.cfg.Auth.TokenIssuer,
			ExpAt: t.Add(ts.cfg.Auth.RefreshTokenLifetime),
			NbfAt: t,
			IssAt: t,
			Payload: token.UserDataClaims{
				UserID: user.ID,
			},
		},
		ts.cfg.Auth.TokenSecretKey,
	)
	if err != nil {
		logger.Error("failed to generate refresh token", "err", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	logger.Info("refresh token generated", "refreshToken", refreshToken)
	return &UserTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
