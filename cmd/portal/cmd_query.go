package main

import (
	"context"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"io"
	"net/url"
	"os"
	"strings"
)

func (a *Application) queryApp(ctx context.Context, query string) (err error) {
	log := plog.Get(ctx)
	log.Println("running query", query)

	target := ""
	target, query = splitQuery(query)
	query = parseQuery(query)
	conn, err := a.Apphost.Query(target, query, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		cancel()
	}()
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		cancel()
	}()
	<-ctx.Done()
	return
}

func splitQuery(targetQuery string) (target string, query string) {
	chunks := strings.SplitN(targetQuery, ":", 2)
	target = chunks[0]
	query = chunks[1]
	return
}

func parseQuery(query string) string {
	if !strings.Contains(query, " -") {
		return query
	}
	vals := url.Values{}
	chunks := strings.Split(query, " ")
	for len(chunks) > 1 {
		if !strings.HasPrefix(chunks[0], "-") {
			break
		}
		vals.Add(chunks[0], chunks[1])
	}
	if len(chunks) > 0 {
		vals["_"] = chunks
	}
	return vals.Encode()
}
