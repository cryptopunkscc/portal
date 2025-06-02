package apphost

type Json[T any] struct {
	Type    string `json:"type"`
	Object  T      `json:"object"`
	Payload []byte `json:"payload"`
}
