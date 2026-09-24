package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func Setup(ch *amqp.Channel) error{

	_, err := ch.QueueDeclare("jobs_queue", true, false, false, false, nil)
	return err
}