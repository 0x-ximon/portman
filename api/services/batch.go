package services

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"unsafe"

	repo "github.com/0x-ximon/portman/api/repositories"
	"github.com/nats-io/nats.go"
)

type BatchService struct {
	client      *nats.Conn
	CreateOrder chan repo.Order
	UpdateOrder chan repo.Order
	CancelOrder chan repo.Order
}

func NewBatchService() (*BatchService, error) {
	queueUrl, ok := os.LookupEnv("QUEUE_URL")
	if !ok {
		return nil, fmt.Errorf("QUEUE_URL is not set")
	}

	client, err := nats.Connect(queueUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid QUEUE_URL: %w", err)
	}

	return &BatchService{
		client:      client,
		CreateOrder: make(chan repo.Order, batchSize),
		UpdateOrder: make(chan repo.Order, batchSize),
		CancelOrder: make(chan repo.Order, batchSize),
	}, nil
}

func (b *BatchService) Start(ctx context.Context) {
	var (
		createBuffer []repo.Order
		updateBuffer []repo.Order
		cancelBuffer []repo.Order
	)

	ticker := time.NewTicker(tickerInterval)
	for {
		select {

		case <-ticker.C:
			b.Send(createBuffer, "CREATE")
			b.Send(updateBuffer, "UPDATE")
			b.Send(cancelBuffer, "CANCEL")

			createBuffer = nil
			updateBuffer = nil
			cancelBuffer = nil

		case order := <-b.CreateOrder:
			createBuffer = append(createBuffer, order)

		case order := <-b.UpdateOrder:
			updateBuffer = append(updateBuffer, order)

		case order := <-b.CancelOrder:
			cancelBuffer = append(cancelBuffer, order)

		case <-ctx.Done():
			ticker.Stop()

			b.Send(createBuffer, "CREATE")
			b.Send(updateBuffer, "UPDATE")
			b.Send(cancelBuffer, "CANCEL")
			return

		}
	}
}

func (b *BatchService) Send(buffer []repo.Order, operation string) error {
	if len(buffer) == 0 {
		return fmt.Errorf("buffer is empty")
	}

	type Headers struct {
		Version     uint8
		Instruction uint8
		Length      uint16
		Nonce       uint32
		Timestamp   uint64
		Source      uint64
		Destination uint64
	}

	type Order struct {
		ID       uint64
		Price    uint64
		Quantity uint64
		Ticker   uint32
		Status   uint8
		Side     uint8
		Mode     uint8
		Flags    uint8
	}

	l := len(buffer) * int(unsafe.Sizeof(Order{}))
	headers := Headers{
		Version:     version,
		Instruction: instruction[operation],
		Length:      uint16(l),
		Nonce:       0,
		Timestamp:   0,
		Source:      0,
		Destination: 0,
	}

	orders := make([]Order, len(buffer))
	for i, o := range buffer {
		// TODO: Properly handle price and quantity conversion to prevent precision loss
		orders[i] = Order{
			ID:       uint64(o.ID),
			Price:    o.Price.BigInt().Uint64(),
			Quantity: o.Quantity.BigInt().Uint64(),
			Ticker:   uint32(o.TickerID),
			Status:   status[o.Status],
			Side:     side[o.Side],
			Mode:     mode[o.Mode],
			Flags:    0x00,
		}
	}

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, headers)
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	err = binary.Write(buf, binary.LittleEndian, orders)
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	err = b.client.Publish("orders.match", buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

const version = 1
const batchSize = 256
const tickerInterval = time.Second * 1

var instruction = map[string]uint8{
	"CREATE": 0,
	"CANCEL": 1,
	"UPDATE": 2,
}

var status = map[repo.OrderStatus]uint8{
	repo.OrderStatusPENDING:   0,
	repo.OrderStatusREJECTED:  1,
	repo.OrderStatusPARTIAL:   2,
	repo.OrderStatusCANCELLED: 3,
	repo.OrderStatusFULFILLED: 4,
}

var mode = map[repo.OrderMode]uint8{
	repo.OrderModeGTC: 0,
	repo.OrderModeFOK: 1,
	repo.OrderModeIOC: 2,
}

var side = map[repo.OrderSide]uint8{
	repo.OrderSideBUY:  0,
	repo.OrderSideSELL: 1,
}
