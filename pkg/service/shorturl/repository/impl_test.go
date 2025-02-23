package repository_test

import (
	"context"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl/repository"
)

type TestSuite struct {
	suite.Suite
	impl         repository.Repository
	dbConnection *sqlx.DB
	pgdb         *embeddedpostgres.EmbeddedPostgres
	driver       database.Driver
	migrate      *migrate.Migrate
}

func (ts *TestSuite) SetupSuite() {
	ts.pgdb = embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Password("password").
		Port(3000).Logger(nil))

	err := ts.pgdb.Start()
	ts.Require().NoError(err)
	ts.dbConnection = sqlx.MustConnect("postgres", "postgres://postgres:password@localhost:3000/postgres?sslmode=disable")
	ts.dbConnection.Exec("CREATE DATABASE dcard;")
	ts.Require().NoError(ts.dbConnection.Close())
	ts.dbConnection = sqlx.MustConnect("postgres", "postgres://postgres:password@localhost:3000/dcard?sslmode=disable")
	ts.driver, err = postgres.WithInstance(ts.dbConnection.DB, &postgres.Config{})
	ts.Require().NoError(err)
	ts.migrate, err = migrate.NewWithDatabaseInstance(
		"file://../../../../deploy/db/migrations",
		"postgres", ts.driver)
	ts.Require().NoError(err)
	ts.impl = repository.New(ts.dbConnection)
}

func (ts *TestSuite) SetupTest() {
	ts.Require().NoError(ts.migrate.Up())
}

func (ts *TestSuite) TearDownTest() {
	ts.Require().NoError(ts.migrate.Down())
}

func (ts *TestSuite) TearDownSuite() {
	err1, err2 := ts.migrate.Close()
	ts.Require().NoError(err1)
	ts.Require().NoError(err2)
	ts.Require().NoError(ts.dbConnection.Close())
	ts.Require().NoError(ts.pgdb.Stop())
}

func (ts *TestSuite) TestCreate() {
	ctx := context.Background()
	testCases := []struct {
		name    string
		short   *domain.ShortURL
		wantErr bool
	}{
		{
			name: "Create short url",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			wantErr: false,
		},
		{
			name: "Create short url with same short code",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with empty short code",
			short: &domain.ShortURL{
				ShortCode:   "",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with empty original url",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with negative expire time",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  -1,
				CreatedTime: 1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with negative created time",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: -1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with zero expire time",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  0,
				CreatedTime: 1,
			},
			wantErr: true,
		},
		{
			name: "Create short url with zero created time",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 0,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		ts.Run(tc.name, func() {
			_, err := ts.impl.Create(ctx, tc.short)
			if tc.wantErr {
				ts.Require().Error(err)
				return
			}
			ts.Require().NoError(err)
		})
	}
}

func (ts *TestSuite) TestGet() {
	ctx := context.Background()
	testCases := []struct {
		name      string
		short     *domain.ShortURL
		expectErr error
	}{
		{
			name: "Get short url",
			short: &domain.ShortURL{
				ShortCode:   "test",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			expectErr: nil,
		},
		{
			name: "Get short url with not exists short code",
			short: &domain.ShortURL{
				ShortCode:   "invalid",
				OriginalURL: "http://test.com",
				ExpireTime:  1,
				CreatedTime: 1,
			},
			expectErr: domain.ErrShortURLNotFound,
		},
	}

	got, err := ts.impl.Create(ctx, testCases[0].short)
	ts.Require().NoError(err)
	ts.Require().NotNil(got)
	for _, tc := range testCases {
		ts.Run(tc.name, func() {
			_, err := ts.impl.Get(ctx, tc.short.ShortCode)
			if tc.expectErr != nil {
				ts.Require().ErrorIs(err, tc.expectErr)
				return
			}
			ts.Require().NoError(err)
		})
	}
}

func TestShortURLSuite(t *testing.T) {
	ts := new(TestSuite)
	suite.Run(t, ts)
}
