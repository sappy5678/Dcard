package shorturl

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"

	"github.com/sappy5678/dcard/pkg/domain"
	"github.com/sappy5678/dcard/pkg/service/shorturl/cache"
	"github.com/sappy5678/dcard/pkg/service/shorturl/repository"
	"github.com/sappy5678/dcard/pkg/service/shorturl/shortcode"
)

type shorturlService struct {
	domain.ShortURLService
	shortcodeGenerator shortcode.Repository
	repo               cache.Repository
	host               string // should get from central config service
	now                func() uint64
}

func New(host string, now func() uint64, shortcodeGenerator shortcode.Repository, repo cache.Repository) domain.ShortURLService {
	return &shorturlService{
		shortcodeGenerator: shortcodeGenerator,
		repo:               repo,
		host:               host,
		now:                now,
	}
}

func Initialize(machineID uint64, host string, db *sqlx.DB, redis rueidis.Client, locker rueidislock.Locker) domain.ShortURLService {
	shortcodeGenerator := shortcode.New(machineID)
	cacheRepo := cache.New(repository.New(db), redis, locker)
	now := func() uint64 {
		now := time.Now().Unix()
		return uint64(now)
	}
	return New(host, now, shortcodeGenerator, cacheRepo)
}
