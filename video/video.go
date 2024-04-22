package video

import (
	"errors"
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/ffmpeg"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAvailableVideos() ([]string, error) {
	var availableVideos []string
	dirEntries, err := os.ReadDir("./")
	if err != nil {
		return nil, fmt.Errorf("video.GetAvailableVideos: could not list directory: %w", err)
	}

	allowedExtensions := map[string]struct{}{
		".mp4":  {},
		".mkv":  {},
		".avi":  {},
		".flv":  {},
		".mov":  {},
		".wmv":  {},
		".ogg":  {},
		".webm": {},
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			fileExtension := filepath.Ext(entry.Name())

			_, ok := allowedExtensions[fileExtension]
			if ok {
				availableVideos = append(availableVideos, entry.Name())
			}
		}
	}

	return availableVideos, nil
}

func GetVideoResolution(videoName string) (string, error) {
	cmd := exec.Command(
		ffmpeg.FfprobePath,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		videoName,
	)
	output, err := common.RunSystemCommand(cmd)
	if err != nil {
		return "", fmt.Errorf("video.GetVideoResolution: ffprobe failed\nstderr: %v\nerr: %w", output, err)
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")

	if output == "" {
		return "", errors.New("video.GetVideoResolution: ffprobe did not return resolution")
	}

	return output, nil
}
