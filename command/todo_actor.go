package command

import (
	"errors"

	"go.uber.org/zap"

	"github.com/lightstaff/go-dddcqrses/common"
	"github.com/lightstaff/go-dddcqrses/events"
	"github.com/lightstaff/go-dddcqrses/messages"
)

// commands
type (
	TodoRegistry struct {
		Message   string
		Completed bool
	}

	TodoMessageChange struct {
		AggregateID string
		Message     string
	}

	TodoComplete struct {
		AggregateID string
		Completed   bool
	}
)

// TodoActor is actor system for todo
type TodoActor struct {
	persistence common.PersistenceContext
	producer    common.MessagingProducerContext
	logger      *zap.SugaredLogger
}

// NewTodoActor is new todo actor
func NewTodoActor(persistence common.PersistenceContext, producer common.MessagingProducerContext, logger *zap.SugaredLogger) *TodoActor {
	return &TodoActor{
		persistence: persistence,
		producer:    producer,
		logger:      logger,
	}
}

// Act is actor action
func (t *TodoActor) Act(command interface{}) error {
	switch command := command.(type) {
	case *TodoRegistry:
		aggregateID, err := common.NewAggregateID()
		if err != nil {
			return err
		}
		entity := NewTodo(aggregateID)
		e, err := events.NewTodoRegistered(aggregateID, command.Message, command.Completed)
		if err != nil {
			return err
		}
		if err := entity.RaiseEvent(e, true); err != nil {
			return err
		}
		if err := t.persistence.Save(entity); err != nil {
			return err
		}
		m, err := messages.NewTodoEventOccurred(entity.AggregateID(), entity.StreamVersion())
		if err != nil {
			return err
		}
		if err := t.producer.Publish(m); err != nil {
			return err
		}
	case *TodoMessageChange:
		entity := NewTodo(command.AggregateID)
		e, err := events.NewTodoMessageChanged(command.AggregateID, command.Message)
		if err != nil {
			return err
		}
		if err := entity.RaiseEvent(e, true); err != nil {
			return err
		}
		if err := t.persistence.Save(entity); err != nil {
			return err
		}
		m, err := messages.NewTodoEventOccurred(entity.AggregateID(), entity.StreamVersion())
		if err != nil {
			return err
		}
		if err := t.producer.Publish(m); err != nil {
			return err
		}
	case *TodoComplete:
		entity := NewTodo(command.AggregateID)
		e, err := events.NewTodoCompleted(command.AggregateID, command.Completed)
		if err != nil {
			return err
		}
		if err := entity.RaiseEvent(e, true); err != nil {
			return err
		}
		if err := t.persistence.Save(entity); err != nil {
			return err
		}
		m, err := messages.NewTodoEventOccurred(entity.AggregateID(), entity.StreamVersion())
		if err != nil {
			return err
		}
		if err := t.producer.Publish(m); err != nil {
			return err
		}
	}

	return errors.New("unknown command")
}
