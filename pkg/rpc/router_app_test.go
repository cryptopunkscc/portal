package rpc

import (
	"context"
	"github.com/cryptopunkscc/astrald/auth/id"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestApp_Run(t *testing.T) {

	// register service
	ctx := context.Background()
	app := NewApp("testApi")
	app.Routes("*")
	app.Func("test", func(i int, b bool) int {
		log.Println("test args", i, b)
		return i
	})
	app.Func("test2", func(s string) string {
		log.Println("test2 args", s)
		return s
	})
	app.Logger(log.New(log.Writer(), "service", 0))
	if err := app.Run(ctx); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Millisecond)

	conn, _ := QueryFlow(id.Identity{}, "testApi")
	//conn := NewRequest(id.Identity{}, "testApi")
	conn.Logger(log.New(log.Writer(), "client ", 0))

	t.Run("Query invalid", func(t *testing.T) {
		err := Command(conn, "asdasdas \n")
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("Query with correct args", func(t *testing.T) {
		r, err := Query[int](conn, "test", 1, true)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, r)
	})

	t.Run("Query with incorrect args", func(t *testing.T) {
		r, err := Query[int](conn, "test asd\n")
		assert.Error(t, err, r)

		r, err = Query[int](conn, "asdasdas\n")
		log.Println(err)
		assert.Error(t, err, r)

		r, err = Query[int](conn, "testasdasdas\n")
		log.Println(err)
		assert.Error(t, err, r)

		r, err = Query[int](conn, "test asd\n")
		assert.Error(t, err, r)

		r, err = Query[int](conn, "testasdasdas\n")
		assert.Error(t, err, r)

		r, err = Query[int](conn, "test asd\n")
		assert.Error(t, err, r)

	})

	t.Run("Query string with newline", func(t *testing.T) {
		r, err := Query[string](conn, "test2", "hello \n world")
		if err != nil {
			t.Fatal(err)
		}
		log.Println(r)
		r, err = Query[string](conn, "test2", "hello \n world")
		if err != nil {
			t.Fatal(err)
		}
		log.Println(r)
	})
}

func TestApp_Run_root(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancel()
	})
	app := NewApp("test")
	app.Logger(log.New(log.Writer(), "service ", 0))
	app.Func("", func(_, identity id.Identity) bool {
		return identity.IsEqual(id.Anyone)
	})
	err := app.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1000000)
	rpc, err := QueryFlow(id.Anyone, "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = rpc.Close()
	})
	rpc.Logger(log.New(log.Writer(), "  client ", 0))

	otherID, _ := id.GenerateIdentity()
	tests := []struct {
		id     id.Identity
		result bool
	}{
		{id.Anyone, true},
		{otherID, false},
	}
	for _, expected := range tests {
		b, err := Query[bool](rpc, "", expected.id)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected.result, b)
	}
	time.Sleep(10000000)
}

func TestApp_Run_subroutine(t *testing.T) {
	ctx := context.Background()
	app := NewApp("test")
	app.Logger(log.New(log.Writer(), "service ", 0))
	app.Func("a", func() (i int, err error) {
		conn, err := QueryFlow(id.Anyone, "test2")
		if err != nil {
			return
		}
		i, err = Query[int](conn, "b")
		return
	})
	err := app.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1000000)
	rpc, err := QueryFlow(id.Anyone, "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = rpc.Close()
	})
	rpc.Logger(log.New(log.Writer(), "  client ", 0))

	app = NewApp("test2")
	app.Logger(log.New(log.Writer(), "service2 ", 0))
	app.Func("b", func() int {
		return 1
	})
	err = app.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1000000)

	i, err := Query[int](rpc, "a")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, i)
}
