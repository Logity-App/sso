package postgres

import (
	"context"
	"fmt"
	"github.com/Logity-App/sso/internal/domain/models"
	"github.com/Logity-App/sso/internal/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	Dialect = `postgres`
)

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg *config.Database) (*Storage, error) {
	if cfg.Enable == false {
		return NewMock(), nil
	}

	pgxCfg, err := pgxpool.ParseConfig(cfg.Dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool cfg parse error: %w", err)
	}
	pgxCfg.MaxConns = int32(cfg.MaxIdleConn)
	pgxCfg.MaxConnLifetime = time.Duration(cfg.MaxLifeTimeConn)

	conn, err := pgxpool.NewWithConfig(context.TODO(), pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool new error: %w", err)
	}

	if err := conn.Ping(context.TODO()); err != nil {
		return nil, fmt.Errorf("pgxpool ping error: %w", err)
	}

	return &Storage{
		db: conn,
	}, nil
}

func NewMock() *Storage {
	return &Storage{}
}

func (s *Storage) Stop() {
	s.db.Close()
}

func (s *Storage) User(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	user.Phone = phone

	return user, nil
}
