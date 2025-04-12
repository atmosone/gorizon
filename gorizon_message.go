package gorizon

import "encoding/json"

type Message struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
}
