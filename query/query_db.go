package query

import "go.uber.org/zap"

// QueryDBContext is query db interface
type QueryDBContext interface {
	FindByID(id string) *TodoQuery
	Save(entity *TodoQuery)
}

// FakeQueryDB is fake query db
type FakeQueryDB struct {
	data   map[string]*TodoQuery
	logger *zap.SugaredLogger
}

// NewFakeQueryDB is new fake query db
func NewFakeQueryDB(logger *zap.SugaredLogger) *FakeQueryDB {
	return &FakeQueryDB{
		data:   make(map[string]*TodoQuery),
		logger: logger,
	}
}

// FindByID is find by id
func (db *FakeQueryDB) FindByID(id string) *TodoQuery {
	if entity, ok := db.data[id]; ok {
		return entity
	}
	return nil
}

// Save is save entity
func (db *FakeQueryDB) Save(entity *TodoQuery) {
	db.data[entity.AggregateID] = entity
}
