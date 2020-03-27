package util

import (
	//"strings"
	"github.com/garyburd/redigo/redis"
	"time"
)

func InitRedisConnection(cfg *Config) (pool *redis.Pool, err error) {
	//redisConfig := cfg.RedisConfig
	pool = &redis.Pool{
		MaxActive:   cfg.RedisConfig.MaxActive,
		MaxIdle:     cfg.RedisConfig.MaxIdle,
		IdleTimeout: time.Duration(cfg.RedisConfig.Timeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.RedisConfig.Host)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	//for key, conf := range redisConfig {
	//
	//	//lowercaseKey := strings.ToLower(key)
	//	host := conf.Host
	//
	//	rds[lowercaseKey] =
	//}

	return
}
