package message

import "encoding/json"

type Type string

const (
	TypeSimple Type = "simple"
)

type Message interface {
	Type() Type
}

type Envelope struct {
	Type Type            `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Simple Message

type SimpleMessage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (m *SimpleMessage) Type() Type {
	return TypeSimple
}

// Ensure all message types implement the Message interface

var _ Message = (*SimpleMessage)(nil)
