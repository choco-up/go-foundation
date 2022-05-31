package databasecache

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Config struct {
	Host       string
	User       string
	Password   string
	DB         int
	DisableTLS bool
	ServerName string
}

func Open(cfg Config) *redis.Client {
	var options = &redis.Options{
		Addr:     cfg.Host,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if !cfg.DisableTLS {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: cfg.ServerName,
		}
	}

	rdb := redis.NewClient(options)

	return rdb
}

func StatusCheck(ctx context.Context, db *redis.Client) error {

	// CHeck if we can ping the database.
	//var pingError error
	for attempts := 1; ; attempts++ {
		val, pingError := db.Ping(ctx).Result()
		if pingError == nil && val == "PONG" {
			break
		}
		fmt.Println(val)
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}
