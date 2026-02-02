package astral_audio_player

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/fs"
	"github.com/cryptopunkscc/astrald/mod/media"
	objects "github.com/cryptopunkscc/astrald/mod/objects/client"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

var Default = &Player{}

func init() {
	Default.ObjectsClient = objects.Default()
}

type Player struct {
	sync.Mutex
	ObjectsClient *objects.Client
	Track         *media.AudioFile
	Location      *fs.FileLocation
	steamer       beep.StreamSeekCloser
	format        beep.Format
}

func (p *Player) SetID(ctx *astral.Context, id *astral.ObjectID) (err error) {
	p.Track = nil
	p.Location = nil
	describeCh, e := p.ObjectsClient.Describe(ctx, id)
	for describe := range describeCh {
		switch d := describe.Descriptor.(type) {
		case *media.AudioFile:
			p.Track = d
		case *fs.FileLocation:
			p.Location = d
		}
	}
	if e != nil && *e != nil {
		return *e
	}
	return nil
}

func (p *Player) SetPath(filePath string) (err error) {
	p.Location = &fs.FileLocation{Path: astral.String16(filePath)}
	return nil
}

func (p *Player) Play(ctx *astral.Context) (err error) {
	var rc io.ReadCloser
	var ext string

	switch {
	case p.Location != nil:
		ext = filepath.Ext(p.Location.Path.String())
		rc, err = os.Open(p.Location.Path.String())
	case p.Track != nil:
		ext = p.Track.Format.String()
		rc, err = p.ObjectsClient.Read(ctx, p.Track.ObjectID, 0, 0)
	default:
		return errors.New("no audio file specified")
	}

	return p.play(rc, ext)
}

func (p *Player) play(rc io.ReadCloser, ext string) (err error) {
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
