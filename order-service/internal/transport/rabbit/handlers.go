package rabbit

import (
	"context"
	"encoding/json"

	"github.com/scmbr/oms/order-service/internal/dto"
	"github.com/scmbr/oms/order-service/internal/models"
	"github.com/streadway/amqp"
)

func (c *Consumer) handle(msg amqp.Delivery) {
	var event dto.SagaEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		msg.Nack(false, false)
		return
	}

	err := c.route(event)
	if err != nil {
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func (c *Consumer) route(event dto.SagaEvent) error {
	switch event.EventType {
	case "inventory.reserved":
		return c.service.UpdateStatus(
			context.Background(),
			event.OrderID,
			models.StatusReserved,
			event.EventID,
		)

	case "inventory.failed":
		return c.service.UpdateStatus(
			context.Background(),
			event.OrderID,
			models.StatusFailed,
			event.EventID,
		)

	case "payment.succeeded":
		return c.service.UpdateStatus(
			context.Background(),
			event.OrderID,
			models.StatusPaid,
			event.EventID,
		)

	case "payment.failed":
		return c.service.UpdateStatus(
			context.Background(),
			event.OrderID,
			models.StatusFailed,
			event.EventID,
		)
	}

	return nil
}
