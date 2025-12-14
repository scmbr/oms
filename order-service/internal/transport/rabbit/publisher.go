package rabbit

import (
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn  *Connection
	queue string
}

func NewPublisher(conn *Connection, queue string) *Publisher {
	return &Publisher{
		conn:  conn,
		queue: queue,
	}
}
func (p *Publisher) Publish(payload []byte, messageID string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.Publish(
		"",
		p.queue,
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
