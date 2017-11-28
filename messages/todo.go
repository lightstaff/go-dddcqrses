package messages

import "github.com/lightstaff/go-dddcqrses/common"

// message types
const MessageTypeTodoEventOccurred = "TodoEventOccurred"

// TodoEventOccurred is todo event occurred message
type TodoEventOccurred struct {
	MessageID     string
	MessageType   string
	AggregateID   string
	StreamVersion int64
}

// NewTodoEventOccurred is new todo event occurred event
func NewTodoEventOccurred(aggregateID string, streamVersion int64) (*TodoEventOccurred, error) {
	messageID, err := common.NewMessageID()
	if err != nil {
		return nil, err
	}
	return &TodoEventOccurred{
		MessageID:     messageID,
		MessageType:   MessageTypeTodoEventOccurred,
		AggregateID:   aggregateID,
		StreamVersion: streamVersion,
	}, nil
}

// GetMessageID is get message id (common.MessageContext interface)
func (m *TodoEventOccurred) GetMessageID() string {
	return m.MessageID
}

// GetMessageType is get message type (common.MessageContext interface)
func (m *TodoEventOccurred) GetMessageType() string {
	return m.MessageType
}
