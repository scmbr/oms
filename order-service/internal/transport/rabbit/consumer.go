package rabbit

import (
	"context"
	"encoding/json"
	"log"

	"github.com/scmbr/oms/order-service/internal/service"
	"github.com/streadway/amqp"
)

type Consumer struct {
	service service.Order
	conn    *Connection
	queue   string
}

func NewConsumer(conn *Connection, queue string, service service.Order) (*Consumer, error) {
	return &Consumer{
		service: service,
		conn:    conn,
		queue:   queue,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		c.queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			go func(d amqp.Delivery) {
				if err := c.handleMessage(d); err != nil {
					log.Printf("failed to handle message: %v", err)
				}
			}(msg)
		}
	}
}

func (c *Consumer) handleMessage(d amqp.Delivery) error {
	type Event struct {
		OrderID string `json:"order_id"`
		Status  string `json:"status"`
		EventID string `json:"event_id"`
	}

	var e Event
	if err := json.Unmarshal(d.Body, &e); err != nil {
		return err
	}

	status, err := c.service.ParseStatus(e.Status)
	if err != nil {
		return err
	}
	return c.service.UpdateStatus(context.Background(), e.OrderID, status, e.EventID)
}
