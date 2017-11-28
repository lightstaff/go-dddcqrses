package common

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// InMemoryMessage is in memory message model
type InMemoryMessage struct {
	Header string
	Data   []byte
}

// MessageContext is message interface
type MessageContext interface {
	GetMessageID() string
	GetMessageType() string
}

// NewMessageID is new event id
func NewMessageID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(err, "IDの生成に失敗しました")
	}

	return id.String(), nil
}

// MessagingProducerContext is messaging producer interface
type MessagingProducerContext interface {
	Publish(m MessageContext) error
}

// FakeMessagingProducer is fake messaging producer
type FakeMessagingProducer struct {
	channel chan<- *InMemoryMessage
	logger  *zap.SugaredLogger
}

// NewFakeMessagingProducer is new fake messaging producer
func NewFakeMessagingProducer(ch chan<- *InMemoryMessage, logger *zap.SugaredLogger) *FakeMessagingProducer {
	return &FakeMessagingProducer{
		channel: ch,
		logger:  logger,
	}
}

// Publish is publish message
func (p *FakeMessagingProducer) Publish(m MessageContext) error {
	d, err := json.Marshal(m)
	if err != nil {
		return err
	}
	msg := &InMemoryMessage{
		Header: m.GetMessageType(),
		Data:   d,
	}
	p.channel <- msg
	p.logger.Infow("publish message", "message", msg)
	return nil
}

// MessagingConsumerContext is messaging consumer interface
type MessagingConsumerContext interface {
	Consume(ctx context.Context, msg chan<- *InMemoryMessage) error
}

// FakeMessagingConsumer is fake massaging consumer
type FakeMessagingConsumer struct {
	channel <-chan *InMemoryMessage
	logger  *zap.SugaredLogger
}

// NewFakeMessagingConsumer is new fake messaging consumer
func NewFakeMessagingConsumer(ch <-chan *InMemoryMessage, logger *zap.SugaredLogger) *FakeMessagingConsumer {
	return &FakeMessagingConsumer{
		channel: ch,
		logger:  logger,
	}
}

// Consume is consume message
func (c *FakeMessagingConsumer) Consume(ctx context.Context, msg chan<- *InMemoryMessage) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-c.channel:
			msg <- m
			c.logger.Infow("consume message", "message", m)
		}
	}
}
