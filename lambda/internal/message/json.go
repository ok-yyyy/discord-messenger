package message

import (
	"encoding/json"
	"fmt"
)

// Unmarshal はJSONデータをMessageインターフェースに変換する。
func Unmarshal(data []byte) (Message, error) {
	var envelope Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message envelope: %w", err)
	}

	switch envelope.Type {
	case TypeSimple:
		var msg SimpleMessage
		if err := json.Unmarshal(envelope.Data, &msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal simple message: %w", err)
		}
		return &msg, nil

	default:
		return nil, fmt.Errorf("unknown message type: %s", envelope.Type)
	}
}

// Marshal はMessageをJSONデータに変換する。
func Marshal(msg Message) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message body: %w", err)
	}

	envelope := Envelope{
		Type: msg.Type(),
		Data: data,
	}

	body, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message envelope: %w", err)
	}

	return body, nil
}
