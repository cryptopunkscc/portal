package stream

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/cryptopunkscc/portal/pkg/plog"
	"github.com/cryptopunkscc/portal/pkg/test"
	"github.com/stretchr/testify/assert"
)

func TestSerializer_Decode(t *testing.T) {
	tests := []struct {
		name       string
		serializer Serializer
		reader     io.Reader
		assert     func(assert.TestingT, any, any, ...any) bool
		actual     any
		expected   any
		wantErr    bool
	}{
		{
			name:     "1",
			actual:   []float64{},
			expected: []float64{1, 2, 3},
			reader:   strings.NewReader("[1,2,3]\n"),
			assert:   assert.ElementsMatch,
		},
		{
			name:     "2",
			actual:   map[string]any{},
			expected: map[string]any{"s": "s", "f": 1.1, "b": true},
			reader:   strings.NewReader(`{"s": "s", "f": 1.1, "b": true}` + "\n"),
			assert:   assert.EqualValues,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serializer := Serializer{}
			serializer.Reader = tt.reader
			serializer.Unmarshal = json.Unmarshal
			if err := serializer.Decode(&tt.actual); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				tt.assert(t, tt.expected, tt.actual)
			}
		})
	}
}

func TestSerializer_Encode(t *testing.T) {
	tests := []struct {
		name       string
		serializer Serializer
		payload    any
		expected   string
		wantErr    bool
	}{
		{
			name:     "1",
			expected: `[1,2,3]` + "\n",
			payload:  []float64{1, 2, 3},
		},
		{
			name:     "2",
			expected: `{"b":true,"f":1.1,"s":"s"}` + "\n",
			payload:  map[string]any{"s": "s", "f": 1.1, "b": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := bytes.NewBuffer(nil)
			serializer := Serializer{
				Writer: buffer,
				Codec: Codec{
					Marshal: json.Marshal,
					Ending:  []byte("\n"),
				},
			}
			if err := serializer.Encode(tt.payload); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.Equal(t, tt.expected, buffer.String())
			}
		})
	}
}

func TestFailure_Json_Marshal_Unmarshal(t *testing.T) {
	e := &Failure{Error: errors.New("==test_error==")}
	b, err := json.Marshal(e)
	test.AssertErr(t, err)
	plog.Println(string(b))
	e = nil
	err = json.Unmarshal(b, &e)
	test.AssertErr(t, err)
}
