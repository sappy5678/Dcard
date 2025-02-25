package cache

import (
	"context"
	"testing"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl/repository"
	redisLocker "github.com/sappy5678/dcard/pkg/utl/locker"
)

type TestSuite struct {
	suite.Suite
	impl       *impl
	redis      rueidis.Client
	locker     rueidislock.Locker
	mockRepo   *repository.MockShortURLRepository
	containers []testcontainers.Container
}

func (ts *TestSuite) SetupSuite() {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis/redis-stack:7.4.0-v3",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForAll(wait.ForLog("Ready to accept connections"), wait.ForListeningPort("6379")),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	ts.Require().NoError(err)

	ts.containers = append(ts.containers, redisC)

	endpoint, err := redisC.Endpoint(ctx, "")
	ts.Require().NoError(err)
	ts.redis, err = rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{endpoint}})
	ts.Require().NoError(err)
	ts.locker, err = redisLocker.New(endpoint)
	ts.Require().NoError(err)
	ts.mockRepo = &repository.MockShortURLRepository{}
	ts.impl = New(ts.mockRepo, ts.redis, ts.locker).(*impl)
}

func (ts *TestSuite) TearDownSuite() {
	testcontainers.CleanupContainer(ts.T(), ts.containers[0])
}

func (ts *TestSuite) SetupTest() {
}

func (ts *TestSuite) TearDownTest() {
	ts.redis.Do(context.Background(), ts.redis.B().Flushall().Build())
}

func (ts *TestSuite) TestCreate() {
	ctx := context.Background()
	ts.mockRepo.CreateFunc = func(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
		return short, nil
	}
	short := &domain.ShortURL{
		ShortCode:   "short",
		OriginalURL: "http://test.com",
	}
	created, err := ts.impl.Create(ctx, short)
	ts.Require().NoError(err)
	ts.Require().Equal(short, created)
	isExist, err := ts.impl.isExist(ctx, short.ShortCode)
	ts.Require().NoError(err)
	ts.Require().True(isExist)
}

func (ts *TestSuite) TestGet_CacheHit() {
	ctx := context.Background()
	shortCode := "exist123"
	expected := &domain.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: "http://exist.com",
	}

	// prefill cache
	err := ts.impl.addBloomFilter(ctx, shortCode)
	ts.Require().NoError(err)
	err = ts.impl.setCache(ctx, expected)
	ts.Require().NoError(err)

	// Mock should not be called
	ts.mockRepo.GetFunc = func(ctx context.Context, code string) (*domain.ShortURL, error) {
		ts.Fail("should not be called")
		return nil, nil
	}

	result, err := ts.impl.Get(ctx, shortCode)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, result)
}

func (ts *TestSuite) TestGet_CacheMiss_DBHit() {
	ctx := context.Background()
	shortCode := "dbExist123"
	expected := &domain.ShortURL{
		ShortCode:   shortCode,
		OriginalURL: "http://db.com",
	}

	// setup Mock and cache
	ts.mockRepo.GetFunc = func(ctx context.Context, code string) (*domain.ShortURL, error) {
		ts.Require().Equal(shortCode, code)
		return expected, nil
	}
	ts.impl.addBloomFilter(ctx, shortCode)

	// validate cache is empty
	cacheKey := ts.impl.getCacheKey(shortCode)
	_, err := ts.redis.Do(ctx, ts.redis.B().Get().Key(cacheKey).Build()).ToString()
	ts.Require().Error(err)

	result, err := ts.impl.Get(ctx, shortCode)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, result)

	// validate cache is filled
	cached, err := ts.impl.getCache(ctx, shortCode)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, cached)
}

func (ts *TestSuite) TestGet_InvalidShortCode() {
	ctx := context.Background()
	shortCode := "invalid"

	ts.mockRepo.GetFunc = func(ctx context.Context, code string) (*domain.ShortURL, error) {
		ts.Fail("should not be called")
		return nil, nil
	}

	_, err := ts.impl.Get(ctx, shortCode)
	ts.Require().ErrorIs(err, domain.ErrShortURLNotFound)

	// validate bloom filter
	isExist, err := ts.impl.isExist(ctx, shortCode)
	ts.Require().NoError(err)
	ts.Require().False(isExist)
}

func (ts *TestSuite) TestGet_InvalidShortCodeFromDB() {
	ctx := context.Background()
	shortCode := "invalid"

	ts.mockRepo.GetFunc = func(ctx context.Context, code string) (*domain.ShortURL, error) {
		return nil, domain.ErrShortURLNotFound
	}
	ts.impl.addBloomFilter(ctx, shortCode) // bloom filter pass invalid short code

	_, err := ts.impl.Get(ctx, shortCode)
	ts.Require().ErrorIs(err, domain.ErrShortURLNotFound)

	// validate bloom filter
	isExist, err := ts.impl.isExist(ctx, shortCode)
	ts.Require().NoError(err)
	ts.Require().True(isExist)
}

func TestCacheSuite(t *testing.T) {
	ts := new(TestSuite)
	suite.Run(t, ts)
}
