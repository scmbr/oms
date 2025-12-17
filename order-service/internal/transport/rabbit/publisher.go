package rabbit

import (
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn     *Connection
	exchange string
}

func NewPublisher(conn *Connection, exchange string) *Publisher {
	return &Publisher{
		conn:     conn,
		exchange: exchange,
	}
}

func (p *Publisher) Publish(routingKey string, payload []byte, messageID string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
			MessageId:   messageID,
			Timestamp:   time.Now(),
		},
	)
}
