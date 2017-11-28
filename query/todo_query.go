package query

import (
	"github.com/lightstaff/go-dddcqrses/common"
	"github.com/lightstaff/go-dddcqrses/events"
)

// TodoQuery is todo query model
type TodoQuery struct {
	AggregateID   string
	Message       string
	Completed     bool
	StreamVersion int64
}

// ApplyEvent is apply event
func (t *TodoQuery) ApplyEvent(e common.EventContext) error {
	switch e := e.(type) {
	case *events.TodoRegistered:
		t.Message = e.Message
		t.Completed = e.Completed
	case *events.TodoMessageChanged:
		t.Message = e.Message
	case *events.TodoCompleted:
		t.Completed = e.Completed
	}
	return nil
}
