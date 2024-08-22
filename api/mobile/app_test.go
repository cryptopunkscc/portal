package mobile

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func Test_callbackReader_Read(t *testing.T) {
	c := make(chan int, 3)
	tc := testCallback{c}
	ar := AsyncReader{testReader{c}}
	ar.Read([]byte{}, tc)
	time.Sleep(1 * time.Millisecond)
	c <- 1
	time.Sleep(3 * time.Millisecond)
	close(c)
	actual := ""
	for i := range c {
		actual += strconv.Itoa(i)
	}
	assert.Equal(t, "123", actual)
}

type testCallback struct{ c chan<- int }
type testReader struct{ c chan<- int }

func (tc testCallback) Result(int, error) { tc.c <- 3 }
func (tr testReader) Read(p []byte) (n int, err error) {
	time.Sleep(2 * time.Millisecond)
	tr.c <- 2
	return 0, nil
}
func (tr testReader) ReadN(int) (arr []byte, err error) { panic("no op") }
