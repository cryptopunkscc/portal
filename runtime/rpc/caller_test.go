package rpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func testHandle(s string) func(a int, b bool) (s string, err error) {
	return func(a int, b bool) (string, error) {
		return fmt.Sprint(s, a, b), nil
	}
}

func TestExecutor_Call(t *testing.T) {
	type fields struct {
		env      []any
		function any
	}
	type args struct {
		args string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult []any
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "function with json positional arguments",
			fields: fields{
				function: testFunc,
			},
			args: args{
				args: `[1, true, "a"]`,
			},
			wantResult: []any{1, true, "a"},
			wantErr:    assert.NoError,
		},
		{
			name: "method with json positional arguments",
			fields: fields{
				function: testStruct{0.1}.method,
			},
			args: args{
				args: `[1, true, "a"]`,
			},
			wantResult: []any{1, true, "a", 0.1},
			wantErr:    assert.NoError,
		},
		{
			name: "method with cli positional arguments",
			fields: fields{
				function: testStruct{0.1}.method,
			},
			args: args{
				args: " 1 true a\n",
			},
			wantResult: []any{1, true, "a", 0.1},
			wantErr:    assert.NoError,
		},
		{
			name: "method with cli mixed arguments ",
			fields: fields{
				function: testFunc7,
			},
			args: args{
				args: " 1 true -s a\n",
			},
			wantResult: []any{1, true, "a"},
			wantErr:    assert.NoError,
		},
		{
			name: "function with json object argument",
			fields: fields{
				function: testFunc2,
			},
			args: args{
				args: `{"i":1}`,
			},
			wantResult: []any{TestArg{I: 1, B: false}},
			wantErr:    assert.NoError,
		},
		{
			name: "function with context and positional argument",
			fields: fields{
				env:      []any{1, context.Background(), true, "ctx"},
				function: testFunc3,
			},
			args: args{
				args: `[1]`,
			},
			wantResult: []any{context.Background(), "ctx", 1},
			wantErr:    assert.NoError,
		},
		{
			name: "function with cli named arguments",
			fields: fields{
				function: testFunc2,
			},
			args: args{
				args: "$ -i 1 -b true -s a \n",
			},
			wantResult: []any{TestArg{I: 1, B: true, TestArg2: TestArg2{S: "a"}}},
			wantErr:    assert.NoError,
		},
		{
			name: "function with many cli named arguments",
			fields: fields{
				function: testFunc4,
			},
			args: args{
				args: " -i 1 -b true -c 1 \n",
			},
			wantResult: []any{TestArg{I: 1, B: true}, TestArg3{C: 1}},
			wantErr:    assert.NoError,
		},
		{
			name: "function with cli positional arguments",
			fields: fields{
				function: testFunc5,
			},
			args: args{
				args: "$ 1 true \n",
			},
			wantResult: []any{TestArgPos{I: 1, B: true}},
			wantErr:    assert.NoError,
		},
		{
			name: "function with many cli mixed arguments",
			fields: fields{
				function: testFunc6,
			},
			args: args{
				args: `$ -s "a" 1 true` + "\n",
			},
			wantResult: []any{TestArgPos{I: 1, B: true}, TestArg2{"a"}},
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewCaller("").With(tt.fields.env...)
			exec.Func(tt.fields.function)
			gotResult, err := exec.Call(NewByteScannerReader(strings.NewReader(tt.args.args)))
			if !tt.wantErr(t, err, fmt.Sprintf("Call(%v)", tt.args.args)) {
				return
			}
			assert.Equalf(t, tt.wantResult, gotResult, "Call(%v)", tt.args.args)
		})
	}
}

func TestExecutor_Call_nested(t *testing.T) {
	caller := NewCaller("").With("s ", 1, true).Func(testHandle)
	r, err := caller.Call(nil)
	if err != nil {
		panic(err)
	}
	log.Print(r)
}

type testStruct struct{ f float64 }

func (ts testStruct) method(i int, b bool, s string) (int, bool, string, float64) {
	return i, b, s, ts.f
}

type TestArg struct {
	I int  `json:"i" name:"i"`
	B bool `json:"b" name:"b"`
	TestArg2
}

type TestArg2 struct {
	S string `json:"s" name:"s"`
}

type TestArg3 struct {
	C int `json:"c" name:"c"`
}

type TestArgPos struct {
	I int  `pos:"1"`
	B bool `pos:"2"`
	*TestArg2
}

func testFunc(i int, b bool, s string) (int, bool, string)            { return i, b, s }
func testFunc2(arg TestArg) TestArg                                   { return arg }
func testFunc3(c context.Context, s string, i int) (any, any, int)    { return c, s, i }
func testFunc4(arg1 TestArg, arg2 TestArg3) (TestArg, TestArg3)       { return arg1, arg2 }
func testFunc5(arg1 TestArgPos) TestArgPos                            { return arg1 }
func testFunc6(arg1 TestArgPos, arg2 TestArg2) (TestArgPos, TestArg2) { return arg1, arg2 }
func testFunc7(i int, b bool, s TestArg2) (int, bool, string)         { return i, b, s.S }
