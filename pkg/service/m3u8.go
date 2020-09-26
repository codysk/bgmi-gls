package service

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"moe.two.bgmi-gls/pkg/common/command"
	"os"
	"os/exec"
	"strings"
	"time"
)

func M3u8VideoIndex(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("%s Not Exist", path)
	}

	totalDuration, err := getVideoDuration(path)
	if err != nil {
		return "", err
	}

	indexFrames, err := getKeyframe(path)
	if err != nil {
		return "", err
	}

	indexFrames = append(indexFrames, totalDuration)
	fragments := calculateFragment(indexFrames)

	targetDuration := largestFragment(fragments).duration

	body := `#EXTM3U
#EXT-X-PLAYLIST-TYPE:VOD
`
	body += fmt.Sprintf("#EXT-X-TARGETDURATION:%v\n", targetDuration.Seconds())

	for _, frag := range fragments {
		body += fmt.Sprintf(
			"#EXTINF:%f,\n%s\n",
			frag.duration.Seconds(),
			fmt.Sprintf(
				"/ts/%s?startAt=%f&duration=%f",
				path,
				frag.start.Seconds(),
				frag.duration.Seconds(),
			),
		)
	}
	body += "#EXT-X-ENDLIST"
	return body, nil
}

func calculateFragment(indexFrames []time.Duration) []VideoFragment {
	indexFramesCount := len(indexFrames) - 1

	fragments := make([]VideoFragment, 0)
	start := indexFrames[0]
	duration := indexFrames[1] - indexFrames[0]
	for i := 1; i < indexFramesCount; {
		if duration >= 10 * time.Second {
			fragments = append(fragments, VideoFragment{
				start:    start,
				duration: duration,
			})
			start = indexFrames[i]
			duration = 0
		}
		duration += indexFrames[i + 1] - indexFrames[i]
		i++
	}
	//last fragment
	fragments = append(fragments, VideoFragment{
		start:    start,
		duration: duration,
	})
	return fragments
}

func largestFragment(fragments []VideoFragment) VideoFragment {
	max := fragments[0]
	for _, frag := range fragments{
		if frag.duration > max.duration {
			max = frag
		}
	}
	return max
}

func getVideoDuration(path string) (time.Duration, error) {
	errWriter := bytes.NewBufferString("")
	outWriter := bytes.NewBufferString("")
	cmd, err := command.NewCommand(
		exec.Command(
			"ffmpeg/ffprobe",
			"-i", path,
			"-v", "error",
			"-show_entries", "format=duration",
			"-of", "default=noprint_wrappers=1:nokey=1",
		),
		nil,
		outWriter,
		errWriter,
	)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()
	_ = cmd.Wait()

	if errWriter.String() != "" {
		log.Info(errWriter)
	}

	outStr := strings.Trim(outWriter.String(), "\n\r\t")
	duration, err := time.ParseDuration(outStr + "s")
	if err != nil {
		return 0, fmt.Errorf("parse time failed: %v", err)
	}

	return duration, nil
}

func getKeyframe(path string) ([]time.Duration, error){
	errWriter := bytes.NewBufferString("")
	outWriter := bytes.NewBufferString("")
	cmd, err := command.NewCommand(
		exec.Command(
			"ffmpeg/ffprobe",
			"-i", path,
			"-v", "error",
			"-skip_frame", "nokey",
			"-select_streams", "v:0",
			"-show_entries", "frame=pkt_pts_time",
			"-of", "csv=print_section=0",
		),
		nil,
		outWriter,
		errWriter,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()
	_ = cmd.Wait()
	if errWriter.String() != "" {
		log.Info(errWriter)
	}

	rawOutputStr := outWriter.String()
	lines := strings.Split(rawOutputStr, "\n")
	lineCount := len(lines)

	keyIndex := make([]time.Duration, 0)
	for i:=0; i < lineCount; i++ {
		line := strings.Trim(lines[i], "\n\r\t")
		if line == "" {
			continue
		}
		duration, err := time.ParseDuration(line + "s")
		if err != nil {
			return nil, err
		}
		keyIndex = append(keyIndex, duration)
	}
	return keyIndex, nil
}
