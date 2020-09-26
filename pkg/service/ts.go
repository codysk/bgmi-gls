package service

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"moe.two.bgmi-gls/pkg/common/command"
	"net/http"
	"os"
	"os/exec"
)

func TSFragment(path string, startAt string, duration string, w http.ResponseWriter) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s Not Exist", path)
	}

	errWriter := bytes.NewBufferString("")
	outReader, outWriter := io.Pipe()

	cmd, err := command.NewCommand(
		exec.Command(
			"ffmpeg/ffmpeg",
			"-v", "error",
			"-ss", startAt,
			"-t", duration,
			"-i", path,
			"-c:v", "nvenc_h264",
			"-c:a", "copy",
			"-f", "hls",
			"-bsf:v", "h264_mp4toannexb",
			"-output_ts_offset", startAt,
			"-",
		),
		nil,
		outWriter,
		errWriter,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	w.Header().Add("Content-Type", "video/MP2T")
	go func() {
		fourKBuffer := make([]byte, 4096)
		flusher := w.(http.Flusher)
		for {
			nbyte, err := outReader.Read(fourKBuffer)
			if err != nil {
				log.Error(err)
				flusher = nil
				return
			}
			_, _ = w.Write(fourKBuffer[:nbyte])
			log.Debug("write %d bytes", nbyte)
			flusher.Flush()
		}
	}()

	_ = cmd.Wait()
	if errWriter.String() != "" {
		log.Info(errWriter)
	}

	return nil
}
