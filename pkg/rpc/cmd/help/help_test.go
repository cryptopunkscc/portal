package help

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/rpc/cmd"
	"github.com/stretchr/testify/assert"
)

func Test_injectHelp(t *testing.T) {
	// given
	handler := cmd.Handler{
		Sub: cmd.Handlers{
			{
				Name: "foo",
				Sub: cmd.Handlers{
					{
						Name: "bar",
					},
				},
			},
		},
	}

	// when
	Inject(&handler)

	// then
	// assert root help
	assertHelpFunc(t, handler.Func)
	assertHasHelp(t, handler)

	// assert foo help
	assert.Equal(t, 2, len(handler.Sub))
	assertHasHelp(t, handler.Sub[0])

	// assert bar help
	assert.Equal(t, 2, len(handler.Sub[0].Sub))
	assertHasHelp(t, handler.Sub[0].Sub[0])
}

var helpName = "help h"
var helpType = func() (h Handler) { return }

func assertHasHelp(t *testing.T, handler cmd.Handler) {
	assert.NotEmpty(t, handler.Sub)
	help := handler.Sub[len(handler.Sub)-1]
	assert.Equal(t, helpName, help.Name)
	assertHelpFunc(t, help.Func)
	assert.Empty(t, help.Sub)
}

func assertHelpFunc(t *testing.T, a any) {
	assert.NotNil(t, a)
	assert.IsType(t, helpType, a)
}

func Test_newHelpFunc(t *testing.T) {
	helpFunc := NewFunc(&handler)
	actual := helpFunc().Handler
	assert.Equal(t, handler, actual)
}

var handler = cmd.Handler{
	Func: nil,
	Name: "foo f",
	Desc: "Foo description",
	Params: cmd.Params{
		{
			Name: "i",
			Type: "int",
			Desc: "Integer",
		},
		{
			Name: "b",
			Type: "bool",
			Desc: "Boolean",
		},
		{
			Type: "string",
			Desc: "String",
		},
	},
	Sub: cmd.Handlers{
		{nil, "bar b", "Bar description", nil, nil},
		{nil, "baz", "", nil, nil},
	},
}

func TestHelp_MarshalCLI(t *testing.T) {
	expected := strings.TrimSpace(`
foo f - Foo description

Parameters:

    -i  [int]     - Integer
    -b  [bool]    - Boolean
    $0  [string]  - String

Commands:

    bar b  - Bar description
    baz
`)
	actual := Handler{handler}.MarshalCLI()
	assert.Equal(t, expected, strings.TrimSpace(actual))
	fmt.Printf("====== actual:\n%s\n======\n", actual)
}
