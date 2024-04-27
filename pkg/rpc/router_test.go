package rpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRouter(t *testing.T) {
	r := NewRouter("").With(context.Background())
	f := func(ctx context.Context, arg testRouterStruct) testRouterStruct {
		return arg
	}
	r.Func("test", f)

	expected := []any{testRouterStruct{I: 1, S: "a"}}
	t.Run("Call", func(t *testing.T) {
		r := r.Query(`test{"i":1,"s":"a"}`)
		result, err := r.Call()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, result)
	})
}

type testRouterStruct struct {
	I int    `json:"i" pos:"1"`
	S string `json:"s" pos:"2"`
}

func TestRouter2(t *testing.T) {
	r := NewRouter("")
	f := func(arg testRouterStruct) testRouterStruct {
		return arg
	}
	r.Func("test", f)

	t.Run("Call correct query", func(t *testing.T) {
		result, err := r.Query(`test{"i":1,"s":"a"}`).Call()
		if err != nil {
			t.Fatal(err)
		}
		log.Println(result)
	})

	t.Run("Call invalid query", func(t *testing.T) {
		result, err := r.Query("test asd \n").Call()
		log.Println(result)
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("Call invalid query", func(t *testing.T) {
		result, err := r.Query("testasd \n").Call()
		log.Println(result, err)
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("Call invalid query", func(t *testing.T) {
		result, err := r.Query("asdaasd \n").Call()
		log.Println(result, err)
		if err == nil {
			t.Fatal()
		}
	})
}
