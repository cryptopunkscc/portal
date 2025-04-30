package rpc

type JsonObject[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}
