package messages

import (
	"encoding/json"
	"errors"

	"github.com/lightstaff/go-dddcqrses/common"
)

// MessageConveter is message to type
func MessageConveter(messageType string, data []byte) (common.MessageContext, error) {
	switch messageType {
	case MessageTypeTodoEventOccurred:
		var m TodoEventOccurred
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		return &m, nil
	}

	return nil, errors.New("unknown message")
}
