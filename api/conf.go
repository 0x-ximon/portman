package main

import (
	"context"
	"fmt"
	"os"

	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumers struct {
	ordersProcessed jetstream.Consumer
}

type Config struct {
	addr   string
	base   string
	pool   *pgxpool.Pool
	mailer *services.MailService
	cacher *services.CacheService
}

func (c *Config) Load(ctx context.Context) error {
	conn, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	c.pool = conn

	addr, ok := os.LookupEnv("ADDR")
	if !ok || addr == "" {
		addr = "127.0.0.1:3001"
	}
	c.addr = addr

	base, ok := os.LookupEnv("BASE")
	if !ok || base == "" {
		base = "http://127.0.0.1:3001"
	}
	c.base = base
	c.addr = addr

	mailer, err := services.NewMailService()
	if err != nil {
		return fmt.Errorf("mailer setup failed: %w", err)
	}
	c.mailer = mailer

	cacher, err := services.NewCacheService()
	if err != nil {
		return fmt.Errorf("cache setup failed: %w", err)
	}
	c.cacher = cacher

	return nil
}
