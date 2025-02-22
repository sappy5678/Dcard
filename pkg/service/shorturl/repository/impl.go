package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/sappy5678/dcard/pkg/domain"

	"github.com/labstack/echo"
)

// postgresql error code define
// http://www.postgresql.org/docs/9.3/static/errcodes-appendix.html

type ShortURL struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ShortURL {
	return &ShortURL{
		db: db,
	}
}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

const createQuery = `INSERT INTO short_url (short_code, original_url, expire_time, created_time) VALUES ($1, $2, $3, $4)`

func (s *ShortURL) Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
	_, err := s.db.ExecContext(ctx, createQuery, short.ShortCode, short.OriginalURL, short.ExpireTime, short.CreatedTime)
	if err != nil {
		return nil, err
	}

	return short, nil
}

const getQuery = `SELECT short_code, original_url, expire_time, created_time FROM short_url WHERE short_code = $1`

func (s *ShortURL) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	var short domain.ShortURL
	err := s.db.GetContext(ctx, &short, getQuery, shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrShortURLNotFound
		}
		return nil, err
	}

	return &short, nil
}
