package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	common_utils "github.com/kholiqcode/go-common/utils"
)

const (
	maxRetries      = 5
	minRetryBackoff = 300 * time.Millisecond
	maxRetryBackoff = 500 * time.Millisecond
	dialTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 3 * time.Second
	minIdleConns    = 20
	idleTimeout     = 12 * time.Second
)

func NewRedisClient(cfg *common_utils.Config) redis.UniversalClient {

	universalClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           []string{cfg.Redis.Host + cfg.Redis.Port},
		Password:        cfg.Redis.Password, // no password set
		DB:              cfg.Redis.DB,       // use default DB
		MaxRetries:      maxRetries,
		MinRetryBackoff: minRetryBackoff,
		MaxRetryBackoff: maxRetryBackoff,
		DialTimeout:     dialTimeout,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		PoolSize:        cfg.Redis.PoolSize,
		MinIdleConns:    minIdleConns,
		PoolTimeout:     cfg.Redis.PoolTimeout * time.Second,
		IdleTimeout:     idleTimeout,
	})

	return universalClient
}
