package cli

//func TestRunner_Main(t *testing.T) {
//	ctx := context.Background()
//	handlers := map[string]*caller.Func{
//		"inc": caller.New(func(a int) int {
//			return a + 1
//		}),
//	}
//	router := rpc.Router{
//		Unmarshalers: []caller.Unmarshaler{
//			json2.Unmarshaler{},
//			clir.Unmarshaler{},
//		},
//		Registry: registry.New[*caller.Func]().AddAll(handlers),
//	}
//	if err := New(router).Run(ctx); err != nil {
//		panic(err)
//	}
//}
//
//func TestRunner_Run(t *testing.T) {
//	ctx := context.Background()
//	tests := []struct {
//		name    string
//		input   string
//		output  string
//		wantErr bool
//	}{
//		{
//			name:   "cli inc",
//			input:  "inc 1",
//			output: "2\n",
//		},
//		{
//			name:   "json inc",
//			input:  "inc[1]",
//			output: "2\n",
//		},
//		{
//			name:   "json inc 2",
//			input:  "test?[1]",
//			output: "2\n",
//		},
//	}
//	handlers := map[string]*caller.Func{
//		"inc": caller.New(func(a int) int {
//			return a + 1
//		}),
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			// given
//			buffer := bytes.NewBuffer(nil)
//			runner := &Runner{
//				conn: rpc.Serializer{
//					Writer:  buffer,
//					Reader:  strings.NewReader(tt.input),
//					Marshal: json.Marshal,
//				},
//				Router: rpc.Router{
//					Unmarshalers: []caller.Unmarshaler{
//						json2.Unmarshaler{},
//						clir.Unmarshaler{},
//					},
//					Registry: registry.New[*caller.Func]().AddAll(handlers),
//				},
//			}
//
//			// when
//			if err := runner.Run(ctx); (err != nil) != tt.wantErr {
//				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
//			}
//
//			// then
//			actual := buffer.String()
//			assert.Equal(t, tt.output, actual)
//		})
//	}
//}
