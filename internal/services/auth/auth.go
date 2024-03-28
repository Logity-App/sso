package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Logity-App/sso/internal/domain/models"
	"github.com/Logity-App/sso/internal/lib/jwt"
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrProvider UserProvider
	cProvider   CodeProvider
	tokenTTL    time.Duration
}

type UserProvider interface {
	User(ctx context.Context, phone string) (*models.User, error)
	SaveUser(ctx context.Context, dto models.SignUpByPhoneDTO) (*models.User, error)
	SaveNewTempUser(ctx context.Context, phone string, code string) (*models.TempUser, error)
	TempUser(ctx context.Context, phone string) (*models.TempUser, error)
	ConfirmTempUser(ctx context.Context, phone string) error
	DeleteTempUser(ctx context.Context, phone string) error
}

type CodeProvider interface {
	Code(ctx context.Context, user models.User) (*models.Code, error)
	GenerateNewCode(ctx context.Context, user models.User, code string) (*models.Code, error)
}

func New(
	log *slog.Logger,
	userProvider UserProvider,
	codeProvider CodeProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrProvider: userProvider,
		cProvider:   codeProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) VerifyNewPhoneNumber(ctx context.Context, phone string) (*models.TempUser, error) {
	user, err := a.usrProvider.User(ctx, phone)

	if errors.Is(err, sql.ErrNoRows) {

		var tempUser *models.TempUser
		tempUser, err = a.usrProvider.TempUser(ctx, phone)

		if errors.Is(err, sql.ErrNoRows) {
			tempUser, err = a.usrProvider.SaveNewTempUser(ctx, phone, a.generateSmsCode())
		}

		if err != nil {
			return nil, err
		}

		return tempUser, nil

	} else if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, fmt.Errorf("%w", "User already exists")
	}

	return nil, nil
}

func (a *Auth) SendSmsCode(ctx context.Context, phone string, code string) error {
	techUser, err := a.usrProvider.TempUser(ctx, phone)

	if err != nil {
		return err
	}

	if techUser.Code != code {
		return fmt.Errorf("%w", "Wrong code!")
	}

	err = a.usrProvider.ConfirmTempUser(ctx, phone)
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) SignUpByPhone(ctx context.Context, dto models.SignUpByPhoneDTO) (*models.User, error) {
	techUser, err := a.usrProvider.TempUser(ctx, dto.Phone)
	if err != nil {
		return nil, err
	}

	if techUser.IsConfirmed == false {
		return nil, fmt.Errorf("%w", "Phone not  confirmed")
	}

	err = a.usrProvider.DeleteTempUser(ctx, dto.Phone)
	if err != nil {
		return nil, err
	}

	user, err := a.usrProvider.SaveUser(ctx, dto)

	if err != nil {
		return nil, err
	}

	app := models.App{}
	token, err := jwt.NewToken(*user, app, a.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	user.Token.AccessToken = token

	return user, nil
}

func (a *Auth) VerifyPhoneNumber(ctx context.Context, phone string) (*models.User, error) {
	user, err := a.usrProvider.User(ctx, phone)

	if err != nil {
		return nil, err
	}

	code, err := a.cProvider.GenerateNewCode(ctx, *user, a.generateSmsCode())

	if err != nil {
		return nil, err
	}

	user.SmsCode = code

	return user, nil
}

func (a *Auth) SignInByPhone(ctx context.Context, phone string, code string) (*models.User, error) {
	user, err := a.usrProvider.User(ctx, phone)

	if err != nil {
		return nil, err
	}

	c, err := a.cProvider.Code(ctx, *user)

	if err != nil {
		return nil, err
	}

	if c.Value != code {
		return nil, fmt.Errorf("%w", "Wrong code!")
	}

	app := models.App{}
	token, err := jwt.NewToken(*user, app, a.tokenTTL)
	if err != nil {

		return nil, fmt.Errorf("%w", err)
	}

	user.Token.AccessToken = token

	return user, nil
}

func (a *Auth) generateSmsCode() string {
	minCode := 100000
	maxCode := 999999
	randCode := rand.Intn(maxCode-minCode) + minCode

	return strconv.FormatInt(int64(randCode), 10)
}
