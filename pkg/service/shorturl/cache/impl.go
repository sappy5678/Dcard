package cache

import (
	"context"
	"encoding/json"

	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl/repository"
)

type impl struct {
	repo   repository.Repository
	redis  rueidis.Client
	locker rueidislock.Locker
}

const (
	bfKey = "bf:shorturl"
	bfCap = 1e10
	bfErr = 1e-6
)

func New(r repository.Repository, redis rueidis.Client, locker rueidislock.Locker) Repository {
	return &impl{
		repo:   r,
		redis:  redis,
		locker: locker,
	}
}

func (im *impl) getCacheKey(shortCode string) string {
	return "shorturl:" + shortCode
}

func (im *impl) Create(ctx context.Context, short *domain.ShortURL) (*domain.ShortURL, error) {
	err := im.addBloomFilter(ctx, short.ShortCode)
	if err != nil {
		return nil, err
	}
	return im.repo.Create(ctx, short)
}

func (im *impl) Get(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	// Check if the short code exists in the bloom filter
	isExist, err := im.isExist(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, domain.ErrShortURLNotFound
	}

	// Get the short URL from the cache
	short, err := im.getCache(ctx, shortCode)
	if err == nil {
		// Cache hit
		return short, nil
	}
	if err != rueidis.Nil {
		return nil, err
	}

	// get lock to avoid thundering herd
	ctx, cancel, err := im.locker.WithContext(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	defer cancel()
	// Get the short URL from the cache again, so request with secondary lock can get from cache
	short, err = im.getCache(ctx, shortCode)
	if err == nil {
		// Cache hit
		return short, nil
	}
	if err != rueidis.Nil {
		return nil, err
	}
	// Cache miss, get data and set it to cache
	short, err = im.repo.Get(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	im.setCache(ctx, short) // even if an error occurs, it is not necessary to return an error

	return short, nil
}

func (im *impl) addBloomFilter(ctx context.Context, shortCode string) error {
	cmd := im.redis.B().BfInsert().Key(bfKey).Capacity(bfCap).Error(bfErr).Items().Item(shortCode).Build()
	if _, err := im.redis.Do(ctx, cmd).AsIntSlice(); err != nil {
		return err
	}
	return nil
}

func (im *impl) setCache(ctx context.Context, short *domain.ShortURL) error {
	key := im.getCacheKey(short.ShortCode)
	jsonBytes, err := json.Marshal(short)
	if err != nil {
		return err
	}
	cmd := im.redis.B().Set().Key(key).Value(string(jsonBytes)).Build()
	return im.redis.Do(ctx, cmd).Error()
}

func (im *impl) getCache(ctx context.Context, shortCode string) (*domain.ShortURL, error) {
	key := im.getCacheKey(shortCode)
	cmd := im.redis.B().Get().Key(key).Build()
	jsonBytes, err := im.redis.Do(ctx, cmd).AsBytes()
	if err != nil {
		return nil, err
	}
	var short domain.ShortURL
	if err := json.Unmarshal(jsonBytes, &short); err != nil {
		return nil, err
	}
	return &short, nil
}

func (im *impl) isExist(ctx context.Context, shortCode string) (bool, error) {
	cmd := im.redis.B().BfExists().Key(bfKey).Item(shortCode).Build()
	isExist, err := im.redis.Do(ctx, cmd).AsBool()
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}

	return true, nil
}
