package command

import (
	"errors"

	"github.com/lightstaff/go-dddcqrses/common"
	"github.com/lightstaff/go-dddcqrses/events"
)

// Todo is Todo aggregate root
type Todo struct {
	*common.AggregateBase
	message   string
	completed bool
}

// NewTodo is new Todo
func NewTodo(id string) *Todo {
	return &Todo{
		AggregateBase: common.NewAggregateBase(id),
	}
}

// Message is message
func (t *Todo) Message() string {
	return t.message
}

// Completed is complated
func (t *Todo) Completed() bool {
	return t.completed
}

// RaiseEvent is raise event
func (t *Todo) RaiseEvent(e common.EventContext, n bool) error {
	switch e := e.(type) {
	case *events.TodoRegistered:
		t.message = e.Message
		t.completed = e.Completed
	case *events.TodoMessageChanged:
		t.message = e.Message
	case *events.TodoCompleted:
		t.completed = e.Completed
	default:
		return errors.New("unknown event")
	}

	if n {
		t.AppendUncommittedEvent(e)
	}

	return nil
}

// Replay is replay events (common.AggreateContext interface)
func (t *Todo) Replay(storedEvents []*common.StoredEvent) error {
	for _, storedEvent := range storedEvents {
		e, err := events.EventConverter(storedEvent.EventType, storedEvent.Data)
		if err != nil {
			return err
		}

		if err := t.RaiseEvent(e, false); err != nil {
			return err
		}

		t.CommitEvent(e)
		t.SetStreamVersion(storedEvent.StreamVersion)
	}

	return nil
}
