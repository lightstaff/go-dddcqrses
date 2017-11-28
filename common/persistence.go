package common

import (
	"encoding/json"
	"sort"

	"go.uber.org/zap"
)

// InMemoryDB is database
type InMemoryDB struct {
	data   map[string][]*StoredEvent
	logger *zap.SugaredLogger
}

// NewInMemoryDB is new in memory db
func NewInMemoryDB(logger *zap.SugaredLogger) *InMemoryDB {
	return &InMemoryDB{
		data:   make(map[string][]*StoredEvent),
		logger: logger,
	}
}

// GetByID is get stored events by id
func (db *InMemoryDB) GetByID(id string) []*StoredEvent {
	results := make([]*StoredEvent, 0)
	results = append(results, db.data[id]...)
	sort.SliceStable(results, func(i, j int) bool {
		switch {
		case results[i].StreamVersion < results[j].StreamVersion:
			return true
		case results[i].StreamVersion > results[j].StreamVersion:
			return false
		default:
			return results[i].OccurredOn < results[j].OccurredOn
		}
	})
	return results
}

// Save is save stored event
func (db *InMemoryDB) Save(e *StoredEvent) {
	db.data[e.AggregateID] = append(db.data[e.AggregateID], e)
}

// PersistenceContext is persistence interface
type PersistenceContext interface {
	ReplayAggregate(a AggregateContext) error
	Save(a AggregateContext) error
}

// FakePersistence is fake persistence
type FakePersistence struct {
	db     *InMemoryDB
	logger *zap.SugaredLogger
}

// NewFakePersistence is new fake persistence
func NewFakePersistence(db *InMemoryDB, logger *zap.SugaredLogger) *FakePersistence {
	return &FakePersistence{
		db:     db,
		logger: logger,
	}
}

// ReplayAggregate is replay aggregate
func (p *FakePersistence) ReplayAggregate(a AggregateContext) error {
	storedEvents := p.db.GetByID(a.AggregateID())
	if err := a.Replay(storedEvents); err != nil {
		return err
	}
	return nil
}

// Save is save aggregate
func (p *FakePersistence) Save(a AggregateContext) error {
	for _, e := range a.UncommittedEvents() {
		d, err := json.Marshal(e)
		if err != nil {
			return err
		}
		storedEvent := &StoredEvent{
			AggregateID:   a.AggregateID(),
			StreamVersion: a.StreamVersion(),
			OccurredOn:    e.GetOccurredOn(),
			EventType:     e.GetEventType(),
			Data:          d,
		}
		p.db.Save(storedEvent)
		a.CommitEvent(e)
	}
	a.SetStreamVersion(a.StreamVersion() + 1)
	return nil
}

// PersistenceQueryContext is persistence query interface
type PersistenceQueryContext interface {
	QueryEvents(id string, base, limit int64) ([]*StoredEvent, error)
}

// FakePersistenceQuery is fake persistence query
type FakePersistenceQuery struct {
	db     *InMemoryDB
	logger *zap.SugaredLogger
}

// NewFakePersistenceQuery is new fake persistence
func NewFakePersistenceQuery(db *InMemoryDB, logger *zap.SugaredLogger) *FakePersistenceQuery {
	return &FakePersistenceQuery{
		db:     db,
		logger: logger,
	}
}

// QueryEvents is query event by id and stream version
func (p *FakePersistenceQuery) QueryEvents(id string, base, limit int64) ([]*StoredEvent, error) {
	results := make([]*StoredEvent, 0)
	storedEvents := p.db.GetByID(id)
	for _, storedEvent := range storedEvents {
		if storedEvent.StreamVersion >= base && storedEvent.StreamVersion <= limit {
			results = append(results, storedEvent)
		}
	}
	return results, nil
}

// StoredEvent is stored event
type StoredEvent struct {
	AggregateID   string
	StreamVersion int64
	OccurredOn    int64
	EventType     string
	Data          []byte
}
