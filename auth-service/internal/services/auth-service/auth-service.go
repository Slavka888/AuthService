package authservice_service

import (
	"AuthService/internal/domain/models"
	"AuthService/internal/lib/jwt"
	"AuthService/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidApp         = errors.New("invalid app")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	UserByEmail(
		ctx context.Context,
		email string,
	) (models.User, error)

	IsAdmin(ctx context.Context,
		uid int64,
	) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// New creates a new instance of the Auth service.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login authenticates a user and returns a JWT token if the credentials are valid.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.LoginUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.Int("app_id", appID),
	)

	log.Info("logging in user")

	user, err := a.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: user not found: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to retrieve user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: failed to retrieve user: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid password", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: invalid password: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: app not found: %w", op, ErrInvalidApp)
		}

		a.log.Error("failed to retrieve app", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: failed to retrieve app: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: failed to generate token: %w", op, err)
	}

	log.Info("user logged in successfully")

	return token, nil
}

// Register registers a new user with the provided email and password, returning the user's ID if successful.
func (a *Auth) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	uid, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: user already exists: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: failed to save user: %w", op, err)
	}

	log.Info("user registered successfully", slog.Int64("user_id", uid))

	return uid, nil
}

// IsAdmin checks if the user with the given ID has administrative privileges.
func (a *Auth) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", uid),
	)

	log.Info(
		"checking if user is admin",
	)

	isAdmin, err := a.userProvider.IsAdmin(ctx, uid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s: user not found: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to check admin status", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: failed to check admin status: %w", op, err)
	}

	log.Info(
		"checked admin status",
		slog.Bool("is_admin", isAdmin),
	)

	return isAdmin, nil
}
