package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/winnerx0/kron/internal/domain"
)

type Publisher struct {
	ch *amqp.Channel

	queueName string

	confirms chan amqp.Confirmation
}

func NewPublisher(ch *amqp.Channel, queueName string) (*Publisher, error) {

	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		return nil, fmt.Errorf("failed to set confirm mode: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	return &Publisher{
		ch:        ch,
		queueName: queueName,
		confirms:  confirms,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, job domain.Job) error {

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = p.ch.PublishWithContext(ctx, "", p.queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		MessageId:    job.ID,
		Timestamp:    time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to publish job: %w", err)
	}

	select {
	case confirm := <-p.confirms:
		if !confirm.Ack {
			return fmt.Errorf("failed to confirm job")
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	}

}
