package events

import (
	"time"

	"github.com/lightstaff/go-dddcqrses/common"
)

// event types
const (
	EventTypeTodoRegistered     = "TodoRegistered"
	EventTypeTodoMessageChanged = "TodoMessageChanged"
	EventTypeTodoCompleted      = "TodoCompleted"
)

// TodoRegistered is todo registered event
type TodoRegistered struct {
	EventID     string
	EventType   string
	OccurredOn  int64
	AggregateID string
	Message     string
	Completed   bool
}

// NewTodoRegistered is new todo registered
func NewTodoRegistered(aggregateID, message string, completed bool) (*TodoRegistered, error) {
	eventID, err := common.NewEventID()
	if err != nil {
		return nil, err
	}
	return &TodoRegistered{
		EventID:     eventID,
		EventType:   EventTypeTodoRegistered,
		OccurredOn:  time.Now().UnixNano(),
		AggregateID: aggregateID,
		Message:     message,
		Completed:   completed,
	}, nil
}

// GetEventID is get event id (common.EventContext interface)
func (e *TodoRegistered) GetEventID() string {
	return e.EventID
}

// GetEventType is get event type (common.EventContext interface)
func (e *TodoRegistered) GetEventType() string {
	return e.EventType
}

// GetOccurredOn is get occurred on (common.EventContext interface)
func (e *TodoRegistered) GetOccurredOn() int64 {
	return e.OccurredOn
}

// TodoMessageChanged is todo message changed event
type TodoMessageChanged struct {
	EventID     string
	EventType   string
	OccurredOn  int64
	AggregateID string
	Message     string
}

// NewTodoMessageChanged is new todo message changed
func NewTodoMessageChanged(aggregateID, message string) (*TodoMessageChanged, error) {
	eventID, err := common.NewEventID()
	if err != nil {
		return nil, err
	}
	return &TodoMessageChanged{
		EventID:     eventID,
		EventType:   EventTypeTodoRegistered,
		OccurredOn:  time.Now().UnixNano(),
		AggregateID: aggregateID,
		Message:     message,
	}, nil
}

// GetEventID is get event id (common.EventContext interface)
func (e *TodoMessageChanged) GetEventID() string {
	return e.EventID
}

// GetEventType is get event type (common.EventContext interface)
func (e *TodoMessageChanged) GetEventType() string {
	return e.EventType
}

// GetOccurredOn is get occurred on (common.EventContext interface)
func (e *TodoMessageChanged) GetOccurredOn() int64 {
	return e.OccurredOn
}

// TodoCompleted is todo completed event
type TodoCompleted struct {
	EventID     string
	EventType   string
	OccurredOn  int64
	AggregateID string
	Completed   bool
}

// NewTodoCompleted is new todo completed
func NewTodoCompleted(aggregateID string, completed bool) (*TodoCompleted, error) {
	eventID, err := common.NewEventID()
	if err != nil {
		return nil, err
	}
	return &TodoCompleted{
		EventID:     eventID,
		EventType:   EventTypeTodoRegistered,
		OccurredOn:  time.Now().UnixNano(),
		AggregateID: aggregateID,
		Completed:   completed,
	}, nil
}

// GetEventID is get event id (common.EventContext interface)
func (e *TodoCompleted) GetEventID() string {
	return e.EventID
}

// GetEventType is get event type (common.EventContext interface)
func (e *TodoCompleted) GetEventType() string {
	return e.EventType
}

// GetOccurredOn is get occurred on (common.EventContext interface)
func (e *TodoCompleted) GetOccurredOn() int64 {
	return e.OccurredOn
}
