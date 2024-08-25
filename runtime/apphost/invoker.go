package apphost

import (
	"context"
	"fmt"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/cryptopunkscc/portal/api/apphost"
	"github.com/cryptopunkscc/portal/api/target"
	"github.com/cryptopunkscc/portal/pkg/flow"
	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/port"
	"regexp"
	"strings"
	"time"
)

type Invoker struct {
	apphost.Client
	Invoke target.Open
	Log    plog.Logger
	Ctx    context.Context
}

func (i Invoker) Query(identity id.Identity, query string) (conn apphost.Conn, err error) {
	if conn, err = i.Client.Query(identity, query); err == nil {
		return
	}
	if identity != id.Anyone {
		return
	}

	i.Log.Println("invoking app for:", query)
	src := strings.Split(query, "?")[0]
	src = strings.TrimPrefix(src, port.New("").String()) // hacky way to handle (remove) dev prefix FIXME
	if query, err = i.invoke(src); err != nil {
		return
	}

	conn, err = flow.RetryT[apphost.Conn](
		i.Ctx, 8188*time.Millisecond,
		func(ii, n int, d time.Duration) (apphost.Conn, error) {
			i.Log.Printf("retry query: %s - %d/%d attempt %v: retry after %v", query, ii+1, n, err, d)
			return i.Client.Query(identity, query)
		},
	)
	if err == nil {
		i.Log.Println("query succeed", conn.Query())
		return
	}
	return
}

func (i Invoker) invoke(query string) (out string, err error) {
	// resolve app [ name | path | port ]
	index := strings.IndexAny(query, ` ?`)
	app := query[:index]
	// invoke app
	var packages []string
	if packages, err = i.Invoke(i.Ctx, app); err != nil {
		return
	}
	// verify query
	if len(packages) != 1 {
		err = fmt.Errorf("required one target for %s, found: %v", query, packages)
		return
	}
	pkg := packages[0]
	if strings.HasPrefix(query, pkg) {
		// query correct
		return query, nil
	}
	// fixup query
	args := query[index:]
	if len(args) == 0 {
		return pkg, nil
	}
	args = regexp.MustCompile(`^\s*:`).ReplaceAllString(args, ".")
	return pkg + args, nil
}
