package beep

import (
	"errors"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cryptopunkscc/portal/apps/player"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

type Player struct {
	sync.Mutex
	steamer beep.StreamSeekCloser
	format  beep.Format
}

var _ player.Audio = &Player{}

func (p *Player) Play(rc io.ReadCloser, ext string) (err error) {
	var decode func(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error)
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)
	switch ext {
	case "mp3":
		decode = mp3.Decode
	case "ogg":
		decode = vorbis.Decode
	case "wav":
		decode = func(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
			return wav.Decode(rc)
		}
	case "flac":
		decode = func(rc io.ReadCloser) (s beep.StreamSeekCloser, format beep.Format, err error) {
			return flac.Decode(rc)
		}
	}

	streamer, format, err := decode(rc)
	if err != nil {
		log.Fatal(err)
	}

	_ = p.Close()
	p.steamer = streamer
	p.format = format
	if err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/60)); err != nil {
		return
	}
	speaker.Play(streamer)

	return nil
}

func (p *Player) Suspend() error {
	return speaker.Suspend()
}

func (p *Player) Resume() error {
	return speaker.Resume()
}

func (p *Player) CurrentTime() time.Duration {
	return p.format.SampleRate.D(p.steamer.Position())
}

func (p *Player) TotalTime() time.Duration {
	return p.format.SampleRate.D(p.steamer.Len())
}

func (p *Player) Seek(duration time.Duration) (err error) {
	switch {
	case p.steamer == nil:
		return errors.New("no audio file set")
	case p.steamer.Len() == 0:
		return errors.New("cannot seek")
	}
	return p.steamer.Seek(p.format.SampleRate.N(duration))
}

func (p *Player) Move(duration time.Duration) (err error) {
	switch {
	case p.steamer == nil:
		return errors.New("no audio file set")
	case p.steamer.Len() == 0:
		return errors.New("cannot seek")
	}
	position := p.format.SampleRate.N(duration) + p.steamer.Position()
	if position > p.steamer.Len() && p.steamer.Len() > 0 {
		position = p.steamer.Len()
	}
	return p.steamer.Seek(position)
}

func (p *Player) Close() (err error) {
	if p.steamer != nil {
		err = p.steamer.Close()
	}
	speaker.Close()
	return
}
