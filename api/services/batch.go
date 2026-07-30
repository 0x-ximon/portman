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
		CreateOrder: make(chan repo.Order, 0xFF),
		UpdateOrder: make(chan repo.Order, 0xFF),
		CancelOrder: make(chan repo.Order, 0xFF),
	}, nil
}

// TODO: Handle the Returned Error and Implement Retry Logic to Prevent Data Loss
func (b *BatchService) Start(ctx context.Context) {
	var (
		createBuffer [0xFF]repo.Order
		updateBuffer [0xFF]repo.Order
		cancelBuffer [0xFF]repo.Order
	)

	var (
		createCount int
		updateCount int
		cancelCount int
	)

	ticker := time.NewTicker(interval)
	for {
		select {

		case <-ticker.C:
			b.Send(createBuffer, "CREATE", createCount)
			b.Send(updateBuffer, "UPDATE", updateCount)
			b.Send(cancelBuffer, "CANCEL", cancelCount)

			createCount = 0
			updateCount = 0
			cancelCount = 0

		case order := <-b.CreateOrder:
			if createCount < 0xFF {
				createBuffer[createCount] = order
				createCount++
			} else {
				b.Send(createBuffer, "CREATE", createCount)
				createCount = 0
			}

		case order := <-b.UpdateOrder:
			if updateCount < 0xFF {
				updateBuffer[updateCount] = order
				updateCount++
			} else {
				b.Send(updateBuffer, "UPDATE", updateCount)
				updateCount = 0
			}

		case order := <-b.CancelOrder:
			if cancelCount < 0xFF {
				cancelBuffer[cancelCount] = order
				cancelCount++
			} else {
				b.Send(cancelBuffer, "CANCEL", cancelCount)
				cancelCount = 0
			}

		case <-ctx.Done():
			ticker.Stop()

			b.Send(createBuffer, "CREATE", createCount)
			b.Send(updateBuffer, "UPDATE", updateCount)
			b.Send(cancelBuffer, "CANCEL", cancelCount)
			return
		}
	}
}

func (b *BatchService) Send(buffer [0xFF]repo.Order, operation string, count int) error {
	if count == 0 {
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

	headers := Headers{
		Version:     version,
		Instruction: instruction[operation],
		Length:      uint16(count * int(unsafe.Sizeof(Order{}))),
		Nonce:       nonce,
		Timestamp:   uint64(time.Now().UnixNano()),
		Source:      0,
		Destination: 0,
	}

	orders := make([]Order, count)
	for i := range count {
		order := buffer[i]

		// TODO: Prevent Precision Loss During price and quantity conversion
		orders[i] = Order{
			ID:       uint64(order.ID),
			Price:    order.Price.BigInt().Uint64(),
			Quantity: order.Quantity.BigInt().Uint64(),
			Ticker:   uint32(order.TickerID),
			Status:   status[order.Status],
			Side:     side[order.Side],
			Mode:     mode[order.Mode],
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
const interval = time.Second * 10

var nonce = uint32(0)

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
