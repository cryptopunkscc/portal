package exec

import (
	"math"
	"time"
)

func AwaitExp(delay time.Duration) (delays []time.Duration) {
	initial := 2
	sum := time.Duration(0)
	for i := initial; true; i++ {
		t := 1 * time.Millisecond * time.Duration(math.Pow(2, float64(i)))
		sum = sum + t
		if sum > delay {
			t -= sum - delay
			delays = append(delays, t)
			return
		}
		delays = append(delays, t)
	}
	return
}
