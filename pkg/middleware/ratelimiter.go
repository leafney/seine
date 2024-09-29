/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     seine
 * @Date:        2024-08-20 16:03
 * @Description:
 */

package middleware

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	rcache "github.com/leafney/rose-cache"
	"github.com/redis/go-redis/v9"
	"log"
)

type RateLimitConfig struct {
	Enabled           bool
	MaxRequestsPerIP  int
	StoreType         string
	MemoryStoreConfig RateLimitMemoryConfig
	RedisStoreConfig  RateLimitRedisConfig
}

func RateLimitMiddleware(limit int) func(*fiber.Ctx) error {
	limiter := redis_rate.NewLimiter()

	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		key := c.IP() // 使用 IP 作为限速的 key

		res, err := limiter.Allow(ctx, key, redis_rate.PerSecond(limit))
		if err != nil {
			log.Printf("rate limit error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		if res.Allowed == 0 {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
		}

		return c.Next()
	}
}

type RateLimitStorer interface {
	Allow(ctx context.Context, identifier string, maxRequests int) (bool, error)
}

type RateLimitMemoryConfig struct {
	CleanMinutes int64
}

func NewRateLimitMemoryStore(config RateLimitMemoryConfig) RateLimitStorer {
	cache, err := rcache.NewCache(config.CleanMinutes)
	if err != nil {
		panic(err)
	}

	return &RateLimitMemoryStore{
		cache: cache,
	}
}

type RateLimitMemoryStore struct {
	cache *rcache.Cache
}

func (s *RateLimitMemoryStore) Allow(ctx context.Context, identifier string, maxRequests int) (bool, error) {

}

// -------

type RateLimitRedisConfig struct {
	Addr     string
	UserName string
	Password string
	DB       int
}

func NewRateLimitRedisStore(config RateLimitRedisConfig) RateLimitStorer {
	//rdb := rredis.MustNewRedis(config.Addr, &rredis.Option{
	//	Pass: config.Password,
	//	Db:   config.DB,
	//})

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Username: config.UserName,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RateLimitRedisStore{
		limiter: redis_rate.NewLimiter(rdb),
	}
}

type RateLimitRedisStore struct {
	limiter *redis_rate.Limiter
}

func (s *RateLimitRedisStore) Allow(ctx context.Context, identifier string, maxRequest int) (bool, error) {

	if maxRequest <= 0 {
		return true, nil
	}

	res, err := s.limiter.Allow(ctx, identifier, redis_rate.PerSecond(maxRequest))
	if err != nil {
		return false, err
	}

	return res.Allowed > 0, nil
}
