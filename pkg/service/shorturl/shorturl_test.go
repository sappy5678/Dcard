package shorturl_test

import (
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl"
	"github.com/sappy5678/dcard/pkg/service/shorturl/cache"
	"github.com/sappy5678/dcard/pkg/service/shorturl/shortcode"
)

var mockHost = "http://test:test"

type TestSuite struct {
	suite.Suite
	impl               domain.ShortURLService
	repo               *cache.MockShortURLCacheRepository
	shortCodeGenerator *shortcode.MockShortCodeIDRepository
	mockNowFn          func() uint64
	mockNow            *time.Time
}

func (ts *TestSuite) SetupSuite() {
	ts.shortCodeGenerator = &shortcode.MockShortCodeIDRepository{}
	ts.repo = &cache.MockShortURLCacheRepository{}
	ts.mockNowFn = func() uint64 {
		return uint64(ts.mockNow.Unix())
	}
	ts.impl = shorturl.New(mockHost, ts.mockNowFn, ts.shortCodeGenerator, ts.repo)
}

func (ts *TestSuite) TearDownSuite() {}

func (ts *TestSuite) TestCreate_NormalCase() {
	now := time.Now()
	ts.mockNow = &now
	expireTime := uint64(ts.mockNow.Add(time.Hour).Unix())

	ts.shortCodeGenerator.NextIDFunc = func() string { return "abc123" }
	expectedShort := &domain.ShortURL{
		ShortCode:   "abc123",
		ShortURL:    mockHost + "/abc123",
		OriginalURL: "https://example.com",
		ExpireTime:  expireTime,
		CreatedTime: uint64(ts.mockNow.Unix()),
	}
	ts.repo.CreateFunc = func(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
		return short, nil
	}

	result, err := ts.impl.Create(context.Background(), "https://example.com", expireTime)

	ts.Require().NoError(err)
	ts.Require().Equal(expectedShort, result)
}

func (ts *TestSuite) TestCreate_ShortCodeGenerationFailure() {
	ts.shortCodeGenerator.NextIDFunc = func() string { return "" }

	_, err := ts.impl.Create(context.Background(), "https://example.com", 0)
	ts.Require().ErrorIs(err, domain.ErrShortURLInvalid)
}

func (ts *TestSuite) TestCreate_RepositoryFailure() {
	ts.shortCodeGenerator.NextIDFunc = func() string { return "def456" }
	ts.repo.CreateFunc = func(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
		return nil, fmt.Errorf("database error")
	}

	_, err := ts.impl.Create(context.Background(), "https://invalid.url", 0)
	ts.ErrorContains(err, "database error")
}

func (ts *TestSuite) TestGet_NormalCase() {
	now := time.Now()
	ts.mockNow = &now
	expireTime := uint64(ts.mockNow.Add(time.Hour).Unix())

	mockShortCode := "mockShortCode"
	ts.repo.GetFunc = func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return &domain.ShortURL{
			ShortCode:   shortCode,
			OriginalURL: "https://example.com",
			ExpireTime:  expireTime,
			CreatedTime: uint64(ts.mockNow.Unix()),
		}, nil
	}
	expect := domain.ShortURL{
		ShortCode:   mockShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", mockHost, mockShortCode),
		OriginalURL: "https://example.com",
		ExpireTime:  expireTime,
		CreatedTime: uint64(ts.mockNow.Unix()),
	}

	result, err := ts.impl.Get(context.Background(), "abc123")
	ts.Require().Equal(expect, result)

	ts.Require().NoError(err)
	ts.Require().NotNil(result)
}

func (ts *TestSuite) TestGet_NotFound() {
	ts.repo.GetFunc = func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return nil, domain.ErrShortURLNotFound
	}

	_, err := ts.impl.Get(context.Background(), "notfound")

	ts.Require().ErrorIs(err, domain.ErrShortURLNotFound)
}

func (ts *TestSuite) TestGet_Expired() {
	now := time.Now()
	ts.mockNow = &now
	expireTime := uint64(ts.mockNow.Add(-time.Hour).Unix())

	ts.repo.GetFunc = func(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
		return &domain.ShortURL{
			ShortCode:   shortCode,
			ShortURL:    fmt.Sprintf("http://test:test/%s", shortCode),
			OriginalURL: "https://example.com",
			ExpireTime:  expireTime,
			CreatedTime: uint64(ts.mockNow.Unix()),
		}, nil
	}

	_, err := ts.impl.Get(context.Background(), "expired")

	ts.Require().ErrorIs(err, domain.ErrShortURLNotFound)
}
