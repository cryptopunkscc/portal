package rpc

import "strings"

func port(chunk string, chunks ...string) string {
	return strings.Join(append([]string{chunk}, chunks...), ".")
}
