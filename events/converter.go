package events

import (
	"encoding/json"
	"errors"

	"github.com/lightstaff/go-dddcqrses/common"
)

// EventConverter is event to type
func EventConverter(eventType string, data []byte) (common.EventContext, error) {
	switch eventType {
	case EventTypeTodoRegistered:
		var e TodoRegistered
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case EventTypeTodoMessageChanged:
		var e TodoMessageChanged
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	case EventTypeTodoCompleted:
		var e TodoCompleted
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}
		return &e, nil
	}
	return nil, errors.New("unknown event")
}
