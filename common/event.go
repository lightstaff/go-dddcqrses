package common

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// EventContext is event interface
type EventContext interface {
	GetEventID() string
	GetEventType() string
	GetOccurredOn() int64
}

// NewEventID is new event id
func NewEventID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(err, "IDの生成に失敗しました")
	}

	return id.String(), nil
}

