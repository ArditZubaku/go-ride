package contracts

import "encoding/json"

// WSMessage is the message structure for the WebSocket.
type WSMessage[T comparable] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type WSDriverMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
