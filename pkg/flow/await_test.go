package flow

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestRetryChan(t *testing.T) {
	c := Await{
		UpTo:  3 * time.Second,
		Delay: 1 * time.Second,
		Mod:   8,
	}.Chan()
	for in := range c {
		log.Println(in)
	}
}

func TestAwaitExp2(t *testing.T) {
	delays := Await{
		UpTo:  5 * time.Second,
		Delay: 50 * time.Millisecond,
		Mod:   8,
	}.exp()
	log.Println(delays)
}

func TestAwaitExp(t *testing.T) {
	tests := []struct {
		name       string
		await      Await
		wantDelays []time.Duration
	}{
		{
			name: "mod 1",
			await: Await{
				UpTo: 50,
				Unit: time.Nanosecond,
				Mod:  1,
			},
			wantDelays: []time.Duration{4, 9, 16},
		},
		{
			name: "mod 2",
			await: Await{
				UpTo: 100,
				Unit: time.Nanosecond,
				Mod:  2,
			},
			wantDelays: []time.Duration{4, 16, 36},
		},
		{
			name: "mod 3",
			await: Await{
				UpTo: 200,
				Unit: time.Nanosecond,
				Mod:  3,
			},
			wantDelays: []time.Duration{9, 36, 81},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDelays := tt.await.exp(); !reflect.DeepEqual(gotDelays, tt.wantDelays) {
				t.Errorf("AwaitExp() = %v, want %v", gotDelays, tt.wantDelays)
			}
		})
	}
}
