package rpc

type JsonObject[T any] struct {
	Type   string `json:"type"`
	Object T      `json:"object"`
}
