package vlc

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"time"

	vlc "github.com/adrg/libvlc-go/v3"
	"github.com/cryptopunkscc/portal/apps/player"
)

type Player struct {
	player     *vlc.Player
	media      *vlc.Media
	mu         sync.Mutex
	closed     bool
	totalDur   time.Duration
	readCloser io.ReadCloser
}

var _ player.Player = &Player{}

func NewPlayer() (*Player, error) {
	if err := vlc.Init("--quiet"); err != nil {
		return nil, err
	}

	p, err := vlc.NewPlayer()
	if err != nil {
		return nil, err
	}

	return &Player{
		player: p,
	}, nil
}

func (p *Player) Play(rc io.ReadCloser, ext string) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("player is closed")
	}

	if p.media != nil {
		p.media.Release()
		p.media = nil
	}
	p.player.Stop()

	rs, ok := rc.(io.ReadSeeker)
	if !ok {
		// TODO use file cache instead of memory cache
		data, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		_ = rc.Close()
		rs = bytes.NewReader(data)
	}
	p.media, err = vlc.NewMediaFromReadSeeker(rs)
	if err != nil {
		return err
	}

	if err = p.player.SetMedia(p.media); err != nil {
		p.media.Release()
		p.media = nil
		return err
	}

	p.readCloser = rc
	p.totalDur = 0

	go func() {
		time.Sleep(800 * time.Millisecond)
		d, _ := p.media.Duration()
		p.mu.Lock()
		p.totalDur = d * time.Millisecond
		p.mu.Unlock()
	}()

	return p.player.Play()
}

func (p *Player) Suspend() error {
	return p.player.SetPause(true)
}

func (p *Player) Resume() error {
	return p.player.SetPause(false)
}

func (p *Player) CurrentTime() time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.player == nil {
		return 0
	}
	t, _ := p.player.MediaTime()
	return time.Duration(t) * time.Millisecond
}

func (p *Player) TotalTime() time.Duration {
	p.mu.Lock()
	d := p.totalDur
	p.mu.Unlock()

	if d == 0 {
		if p.media != nil {
			ms, _ := p.media.Duration()
			d = ms * time.Millisecond
		}
	}
	return d
}

func (p *Player) Seek(d time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.player == nil || p.media == nil {
		return errors.New("no media loaded")
	}
	return p.player.SetMediaTime(int(d / time.Millisecond))
}

func (p *Player) Move(delta time.Duration) error {
	cur := p.CurrentTime()
	return p.Seek(cur + delta)
}

func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}
	p.closed = true

	if p.player != nil {
		p.player.Stop()
		p.player.Release()
	}
	if p.media != nil {
		p.media.Release()
	}
	if p.readCloser != nil {
		_ = p.readCloser.Close()
	}

	vlc.Release()

	return nil
}
