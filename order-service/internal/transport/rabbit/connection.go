package rabbit

import (
	"github.com/streadway/amqp"
)

type Connection struct {
	Conn *amqp.Connection
}

func NewConnection(cfg Config) (*Connection, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(
		cfg.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	ch.Close()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Connection{Conn: conn}, nil
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	return c.Conn.Channel()
}

func (c *Connection) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}
