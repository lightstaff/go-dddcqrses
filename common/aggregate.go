package common

import (
	"sort"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// AggregateContext is aggregate interface
type AggregateContext interface {
	AggregateID() string
	StreamVersion() int64
	SetStreamVersion(int64)
	UncommittedEvents() []EventContext
	AppendUncommittedEvent(EventContext)
	CommitEvent(EventContext)
	CommittedEvents() []EventContext
	Replay([]*StoredEvent) error
}

// NewAggregateID is new aggregate id
func NewAggregateID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(err, "IDの生成に失敗しました")
	}

	return id.String(), nil
}

// AggregateBase is aggregate
type AggregateBase struct {
	aggregateID         string
	streamVersion       int64
	uncommittedEventMap map[string]EventContext
	committedEventMap   map[string]EventContext
}

// NewAggregateBase is new aggregate base
func NewAggregateBase(aggregateID string) *AggregateBase {
	return &AggregateBase{
		aggregateID:         aggregateID,
		streamVersion:       int64(0),
		uncommittedEventMap: make(map[string]EventContext),
		committedEventMap:   make(map[string]EventContext),
	}
}

// AggregateID is aggregate id
func (a *AggregateBase) AggregateID() string {
	return a.aggregateID
}

// StreamVersion is stream version
func (a *AggregateBase) StreamVersion() int64 {
	return a.streamVersion
}

// SetStreamVersion is set stream version
func (a *AggregateBase) SetStreamVersion(value int64) {
	a.streamVersion = value
}

// UncommittedEvents is get uncommitted events
func (a *AggregateBase) UncommittedEvents() []EventContext {
	if a.uncommittedEventMap == nil {
		a.uncommittedEventMap = make(map[string]EventContext)
	}
	results := make([]EventContext, 0, len(a.uncommittedEventMap))
	for _, event := range a.uncommittedEventMap {
		results = append(results, event)
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].GetOccurredOn() < results[j].GetOccurredOn()
	})
	return results
}

// AppendUncommittedEvent is append uncommitted event
func (a *AggregateBase) AppendUncommittedEvent(event EventContext) {
	if a.uncommittedEventMap == nil {
		a.uncommittedEventMap = make(map[string]EventContext)
	}
	a.uncommittedEventMap[event.GetEventID()] = event
}

// CommittedEvents is get committed events
func (a *AggregateBase) CommittedEvents() []EventContext {
	if a.committedEventMap == nil {
		a.committedEventMap = make(map[string]EventContext)
	}
	results := make([]EventContext, 0, len(a.committedEventMap))
	for _, event := range a.committedEventMap {
		results = append(results, event)
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].GetOccurredOn() < results[j].GetOccurredOn()
	})
	return results
}

// CommitEvent is commit event
func (a *AggregateBase) CommitEvent(event EventContext) {
	if a.uncommittedEventMap == nil {
		a.uncommittedEventMap = make(map[string]EventContext)
	}
	delete(a.uncommittedEventMap, event.GetEventID())
	if a.committedEventMap == nil {
		a.committedEventMap = make(map[string]EventContext)
	}
	a.committedEventMap[event.GetEventID()] = event
}
