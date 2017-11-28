package query

import (
	"github.com/lightstaff/go-dddcqrses/common"
	"github.com/lightstaff/go-dddcqrses/events"
	"github.com/lightstaff/go-dddcqrses/messages"
	"go.uber.org/zap"
)

// TodoActor is actor system to todo
type TodoActor struct {
	persistenceQuery common.PersistenceQueryContext
	queryDB          QueryDBContext
	logger           *zap.SugaredLogger
}

// NewTodoActor new todo actor
func NewTodoActor(persistenceQuery common.PersistenceQueryContext, queryDB QueryDBContext, logger *zap.SugaredLogger) *TodoActor {
	return &TodoActor{
		persistenceQuery: persistenceQuery,
		queryDB:          queryDB,
		logger:           logger,
	}
}

// Act is todo actor action
func (t *TodoActor) Act(msg *messages.TodoEventOccurred) error {
	target := t.queryDB.FindByID(msg.AggregateID)
	if target == nil {
		target = &TodoQuery{
			AggregateID: msg.AggregateID,
		}
	}
	storedEvents, err := t.persistenceQuery.QueryEvents(msg.AggregateID, target.StreamVersion, msg.StreamVersion)
	if err != nil {
		return err
	}
	t.logger.Infow("apply events", "events", storedEvents)
	for _, storedEvent := range storedEvents {
		e, err := events.EventConverter(storedEvent.EventType, storedEvent.Data)
		if err != nil {
			return err
		}
		target.ApplyEvent(e)
		target.StreamVersion = storedEvent.StreamVersion
		t.logger.Infow("apply event", "event", e)
	}
	t.queryDB.Save(target)
	return nil
}
