package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Logity-App/sso/internal/domain/models"
	"github.com/Logity-App/sso/internal/pkg/config"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
)

var (
	dialect = goqu.Dialect("postgres")
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Database) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func NewMock() *Storage {
	return &Storage{}
}

func (s *Storage) Stop() {
	s.db.Close()
}

func (s *Storage) User(ctx context.Context, phone string) (*models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, phone, default_tag, birthday_date FROM users WHERE phone = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, phone)

	u := new(models.User)
	err = row.Scan(&u.ID, &u.Phone, &u.DefaultTag, &u.BirthdayDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

const (
	saveUserQuery = `INSERT INTO users (phone, default_tag, birthday_date)
							VALUES ($1, $2, $3) returning id, phone;`
)

func (s *Storage) SaveUser(ctx context.Context, dto models.SignUpByPhoneDTO) (*models.User, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare(saveUserQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.QueryContext(ctx, dto.Phone, dto.DefaultTag, dto.BirthdayDate)
	if err != nil {

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	u := new(models.User)
	for res.Next() {
		if err := res.Scan(&u.ID, &u.Phone); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (s *Storage) TempUser(ctx context.Context, phone string) (*models.TempUser, error) {
	const op = "storage.sqlite.TempUser"

	stmt, err := s.db.Prepare("SELECT id, phone, code, is_confirmed FROM temp_users WHERE phone = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, phone)

	u := new(models.TempUser)
	err = row.Scan(&u.ID, &u.Phone, &u.Code, &u.IsConfirmed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return u, nil
}

func (s *Storage) ConfirmTempUser(ctx context.Context, phone string) error {
	const op = "storage.ConfirmTempUser"

	stmt, err := s.db.Prepare("UPDATE temp_users SET is_confirmed = true WHERE phone = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx, phone); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTempUser(ctx context.Context, phone string) error {
	const op = "storage.sqlite.DeleteTempUser"

	stmt, err := s.db.Prepare("DELETE FROM temp_users WHERE phone = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx, phone); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

const (
	saveNewTempUserQuery = `INSERT INTO temp_users (phone, code)
							VALUES ($1, $2) returning id, phone, code;`
)

func (s *Storage) SaveNewTempUser(ctx context.Context, phone string, code string) (*models.TempUser, error) {
	const op = "storage.SaveNewTempUser"

	stmt, err := s.db.Prepare(saveNewTempUserQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.QueryContext(ctx, phone, code)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	tempUser := new(models.TempUser)
	for res.Next() {
		if err := res.Scan(&tempUser.ID, &tempUser.Phone, &tempUser.Code); err != nil {
			return nil, err
		}
	}

	return tempUser, nil
}

// TODO add  $5::timestamp
const (
	generateNewCodeQuery = `INSERT INTO sms_codes (code, user_id)
							VALUES ($1, $2)
							ON CONFLICT ON CONSTRAINT code_user_id_ukey
							DO UPDATE SET code = EXCLUDED.code;`
)

func (s *Storage) GenerateNewCode(ctx context.Context, user models.User, code string) (*models.Code, error) {
	const op = "storage.GenerateNewCode"

	stmt, err := s.db.Prepare(generateNewCodeQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, code, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &models.Code{
		Value:  code,
		UserID: user.ID,
	}, nil
}

func (s *Storage) Code(ctx context.Context, user models.User) (*models.Code, error) {
	const op = "storage.sqlite.Code"

	stmt, err := s.db.Prepare("SELECT id, code, user_id FROM sms_codes WHERE user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, user.ID)

	c := new(models.Code)
	err = row.Scan(&c.ID, &c.Value, &c.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, "ErrUserNotFound")
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}
