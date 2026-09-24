package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	Ch   *amqp.Channel
	Conn *amqp.Connection
}

func NewRabbitMQClient(url string) *RabbitMQClient {
	conn, err := amqp.Dial(url)

	if err != nil {
		log.Fatal("Error creating connection to rabbitmq ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creating channel ", err)
	}

	return &RabbitMQClient{
		Ch:   channel,
		Conn: conn,
	}
}

func (r *RabbitMQClient) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {

	ch, err := r.Ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (r *RabbitMQClient) PublishMessage(queueName string, body []byte) error {

	err := r.Ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})

	return err
}

func (r *RabbitMQClient) Close() error {

	if err := r.Ch.Close(); err != nil {
		return err
	}

	return r.Conn.Close()
}
