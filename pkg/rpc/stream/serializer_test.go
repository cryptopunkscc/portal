package stream

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"strings"
	"testing"
	"time"
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
			serializer := Serializer{
				Reader:    tt.reader,
				Unmarshal: json.Unmarshal,
			}
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
				Writer:  buffer,
				Marshal: json.Marshal,
				Ending:  []byte("\n"),
			}
			if err := serializer.Encode(tt.payload); (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				assert.Equal(t, tt.expected, buffer.String())
			}
		})
	}
}

func TestTest(t *testing.T) {
	reader, writer := io.Pipe()

	go func() {
		write, err := writer.Write([]byte{})
		log.Println("write", write, err)
	}()

	b := make([]byte, 10)
	read, err := reader.Read(b)
	log.Println("read", read, err)
	time.Sleep(100 * time.Millisecond)
}
