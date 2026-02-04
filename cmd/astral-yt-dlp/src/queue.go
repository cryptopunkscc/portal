package astral_yt_dlp

import (
	"fmt"
	"slices"
	"sync"

	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/api"
)

type Queue struct {
	sync.RWMutex
	sig.Queue[astral_yt_dlp.Request]
	items []astral_yt_dlp.Request
}

func (q *Queue) Push(request astral_yt_dlp.Request) error {
	q.Lock()
	defer q.Unlock()
	if slices.Contains(q.items, request) {
		return fmt.Errorf("already downloading %v", request)
	}
	q.items = append(q.items, request)
	if len(q.items) == 1 {
		q.Queue.Push(q.items[0])
	}
	return nil
}

func (q *Queue) Done(request astral_yt_dlp.Request) bool {
	q.Lock()
	defer q.Unlock()

	l := len(q.items)
	q.items = slices.DeleteFunc(q.items,
		func(r astral_yt_dlp.Request) bool { return r == request },
	)
	if len(q.items) > 0 && len(q.items) != l {
		q.Queue.Push(q.items[0])
	}
	return len(q.items) == 0
}
