package services

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreService struct {
	conn *grpc.ClientConn
}

func NewCoreService() *CoreService {
	return &CoreService{}
}

func (c *CoreService) Connect() (*grpc.ClientConn, error) {
	if c.conn != nil {
		return c.conn, nil
	}

	corePort, ok := os.LookupEnv("CORE_PORT")
	if !ok {
		corePort = "50051"
	}

	coreAddr := net.JoinHostPort(os.Getenv("CORE_HOST"), corePort)
	// TODO: Implement secure connection
	conn, err := grpc.NewClient(coreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to core service: %w", err)
	}

	c.conn = conn
	return conn, err
}
