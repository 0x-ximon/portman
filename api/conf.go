package main

import (
	"context"
	"fmt"
	"os"

	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	addr    string
	pool    *pgxpool.Pool
	web3    *services.Web3Service
	mailer  *services.MailService
	cacher  *services.CacheService
	batcher *services.BatchService
}

func (c *Config) Load(ctx context.Context) error {
	conn, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	c.pool = conn

	addr, ok := os.LookupEnv("ADDR")
	if !ok {
		return fmt.Errorf("ADDR is not set")
	}
	c.addr = addr

	web3, err := services.NewWeb3Service()
	if err != nil {
		return fmt.Errorf("web3 setup failed: %w", err)
	}
	c.web3 = web3

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

	batcher, err := services.NewBatchService()
	if err != nil {
		return fmt.Errorf("batch setup failed: %w", err)
	}
	c.batcher = batcher

	return nil
}
