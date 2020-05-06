package typredis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typcfg"
)

// Redis module
type Redis struct {
	typapp.Provider
	typapp.Destroyer
	typapp.Preparer
	typcfg.Configurer
}

// Module of redis
func Module() *Redis {
	return &Redis{
		Provider:  typapp.NewConstructor("", Connect),
		Destroyer: typapp.NewDestructor(Disconnect),
		Preparer:  typapp.NewPreparation(Ping),
	}
}

// Connect to redis server
func Connect(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:           cfg.Password,
		DB:                 cfg.DB,
		PoolSize:           cfg.PoolSize,
		DialTimeout:        cfg.DialTimeout,
		ReadTimeout:        cfg.ReadWriteTimeout,
		WriteTimeout:       cfg.ReadWriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		MaxConnAge:         cfg.MaxConnAge,
	})
}

// Ping redis server
func Ping(client *redis.Client) (err error) {
	if err = client.Ping().Err(); err != nil {
		return fmt.Errorf("Redis: Ping: %w", err)
	}
	return
}

// Disconnect from service server
func Disconnect(client *redis.Client) (err error) {
	if err = client.Close(); err != nil {
		return fmt.Errorf("Redis: Disconnect: %w", err)
	}
	return
}
