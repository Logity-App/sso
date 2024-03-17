package auth

import (
	"context"
	"github.com/Logity-App/sso/internal/domain/models"
	"log/slog"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	tokenTTL    time.Duration
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrProvider: userProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) VerifyNewPhoneNumber(ctx context.Context, phone string) (*models.Code, error) {
	_, err := a.usrProvider.User(ctx, phone)
	if err != nil {
		return nil, err
	}

	return &models.Code{
		Value:        "4343",
		ExpirationAt: 1111111,
	}, nil
}
