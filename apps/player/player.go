package player

import (
	"io"
	"time"
)

type Audio interface {
	Play(rc io.ReadCloser, ext string) (err error)
	CurrentTime() time.Duration
	TotalTime() time.Duration
	Seek(duration time.Duration) (err error)
	Move(duration time.Duration) (err error)
	Close() (err error)
	Suspend() error
	Resume() error
}

type Video interface {
	Audio
	Fullscreen(bool) (err error)
	IsFullscreen() (bool, error)
}
