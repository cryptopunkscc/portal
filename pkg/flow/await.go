package flow

import (
	"context"
	"time"
)

type Await struct {
	UpTo  time.Duration
	Delay time.Duration
	Unit  time.Duration
	Mod   int
	Ctx   context.Context
}

func (a Await) exp() (delays []time.Duration) {
	if a.Mod == 0 {
		a.Mod = 4
	}
	if a.Unit == 0 {
		a.Unit = time.Millisecond
	}
	if a.Delay > 0 {
		delays = append(delays, a.Delay)
	}
	initial := 2
	sum := time.Duration(0)
	for i := initial; true; i++ {
		if i%a.Mod != 0 {
			continue
		}
		f := float64(i)
		t := time.Duration(f*f) * a.Unit
		sum = sum + t
		if sum > a.UpTo {
			return
		}
		delays = append(delays, t)
	}
	return
}

func (a Await) Chan() <-chan RetryIn {
	if a.Ctx == nil {
		a.Ctx = context.Background()
	}
	delays := a.exp()
	n := len(delays) - 1
	c := make(chan RetryIn, n)
	ctx, cancel := context.WithCancel(a.Ctx)
	go func() {
		<-ctx.Done()
		close(c)
	}()
	go func() {
		time.Sleep(delays[0])
		for _, d := range delays[1:] {
			if ctx.Err() != nil {
				return
			}
			c <- RetryIn{Remains: n, Delay: d}
			time.Sleep(d)
			n--
		}
		if ctx.Err() != nil {
			return
		}
		c <- RetryIn{Remains: n, Delay: 0}
		cancel()
		time.Sleep(10 * time.Millisecond)
		for range c {
		}
	}()
	return c
}

type RetryIn struct {
	Remains int
	Delay   time.Duration
}
