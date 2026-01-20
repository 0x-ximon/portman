package main

import (
	"context"
	"net"
	"os"

	"github.com/0x-ximon/portman/api/services"
	"github.com/jackc/pgx/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/grpc"
)

type Consumers struct {
	ordersProcessed jetstream.Consumer
}

type Config struct {
	addr     string
	dbConn   *pgx.Conn
	natsConn *nats.Conn
	coreConn *grpc.ClientConn
}

func (c *Config) Load(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DB_URL"))
	if err != nil {
		return err
	}
	c.dbConn = conn

	coreConn, err := services.NewCoreService().Connect()
	if err != nil {
		return err
	}
	c.coreConn = coreConn

	natsConn, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return err
	}
	c.natsConn = natsConn

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3001"
	}

	addr := net.JoinHostPort(os.Getenv("HOST"), port)
	c.addr = addr

	return nil
}

func (c *Config) GetConsumers(ctx context.Context) (*Consumers, error) {
	js, err := jetstream.New(c.natsConn)
	if err != nil {
		return nil, err
	}

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "orders",
		Subjects: []string{"orders.*"},
	})
	if err != nil {
		return nil, err
	}

	var consumers Consumers

	{
		cons, err := s.CreateConsumer(ctx, jetstream.ConsumerConfig{
			FilterSubject: "orders.processed",
		})

		if err != nil {
			return nil, err
		}

		consumers.ordersProcessed = cons
	}

	{
		cons, err := s.CreateConsumer(ctx, jetstream.ConsumerConfig{
			FilterSubject: "orders.processed",
		})

		if err != nil {
			return nil, err
		}

		consumers.ordersProcessed = cons
	}

	return &consumers, nil

}
