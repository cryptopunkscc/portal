package astral_yt_dlp

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/sig"
	"github.com/cryptopunkscc/portal/cmd/astral-yt-dlp/api"
	"github.com/cryptopunkscc/portal/pkg/os"
)

type Download struct {
	progress sig.Queue[astral_yt_dlp.Progress]
}

type RunOpts struct {
	Dir   string
	Audio bool
}

func (d *Download) Reset() {
	d.progress.Close()
	d.progress = sig.Queue[astral_yt_dlp.Progress]{}
	return
}

func (d *Download) Run(ctx *astral.Context, request astral_yt_dlp.Request) (err error) {
	progress := astral_yt_dlp.Progress{Status: "started"}
	progress.Request = &request
	d.progress.Push(progress)
	defer func() {
		if err != nil {
			progress.Status = astral.String8(err.Error())
		}
		d.progress.Push(progress)
	}()
	args := []string{
		"--update-to", "stable",
		"--newline",
		"--progress-template", "PROG:%(progress._percent_str)s|%(progress._downloaded_bytes)s|%(progress._total_bytes)s|%(progress._speed_str)s|%(progress._eta_str)s",

		//"--merge-output-format", "mkv",
		"--embed-thumbnail",
		"--embed-metadata",
		"--embed-chapters",
		"--embed-subs",
		"--write-info-json",
	}
	if request.Audio {
		args = append(args, "-x", "bestaudio/best", "--audio-format", "mp3")
	} else {
		args = append(args, "-f", "bestvideo+bestaudio/best")
	}
	if len(request.Dir) > 0 {
		args = append(args, "-P", os.Abs(request.Dir.String()))
	}
	args = append(args, request.Url.String())

	cmd := exec.CommandContext(ctx, "yt-dlp", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "PROG:"):
			if err = progress.UnmarshalText([]byte(line)); err != nil {
				progress.Status = astral.String8(err.Error())
			}
			d.progress.Push(progress)
		case strings.Contains(line, "[download] Destination:"):
			fmt.Println("file:", line)
		}
	}

	if err = cmd.Wait(); err == nil {
		progress.Status = "completed"
	}
	return
}

func (d *Download) Progress(ctx context.Context, follow bool) <-chan astral_yt_dlp.Progress {
	ch := make(chan astral_yt_dlp.Progress)
	go func() {
		defer close(ch)
		for {
			for progress := range sig.Subscribe(ctx, &d.progress) {
				ch <- progress
			}
			if !follow {
				return
			}
		}
	}()
	return ch
}
